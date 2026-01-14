package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pkg/log/model"
	"pkg/mosquitto"
	exporterEndpoint "server/internal/modules/exporter/endpoint"
	"server/internal/utils/errors"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/pressly/goose/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shopspring/decimal"

	"pkg/database/pgsql"
	"pkg/http/router"
	httpServer "pkg/http/server"

	"pkg/log"
	"pkg/migrator"
	"pkg/panicRecover"
	"pkg/trace"
	"server/internal/config"
	exporterRepository "server/internal/modules/exporter/repository"
	exporterService "server/internal/modules/exporter/service"
	pgsqlMigrations "server/migrations/pgsql"
)

const version = "@{version}"
const build = "@{build}"
const buildDate = "@{buildDate}"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// Основной контекст приложения
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Перехватываем возможную панику
	defer func() {
		panicRecover.PanicRecover(func(err error) {
			log.Fatal(err)
		})
	}()

	// Получаем конфиг
	conf := config.Load()

	// Получаем имя хоста
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	// Инициализируем логгер
	if err := log.InitDefaultLogger(
		model.SystemInfo{
			Version:     version,
			Build:       build,
			ServiceName: conf.ServiceName,
			Env:         conf.Environment,
			Hostname:    hostname,
			BuildDate:   buildDate,
		},
		conf.Logger,
	); err != nil {
		return err
	}

	// Инициализируем все синглтоны
	log.Info("Инициализируем синглтоны")
	if err = initSingletons(conf); err != nil {
		return err
	}

	log.Info("Инициализируем трейсер")
	if err = trace.StartTracing(conf.Tracer, conf.ServiceName); err != nil {
		return err
	}

	// Подключаемся к mosquitto
	log.Info("Инициализируем mosquitto")
	mosquitto, err := mosquitto.NewClientMosquitto(ctx, conf.Mosquitto)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}

	// Подключаемся к базе данных
	log.Info("Подключаемся к postgresql")
	pgsql, err := pgsql.NewClientPgsql(ctx, conf.Pgsql)
	if err != nil {
		return err
	}
	defer pgsql.Close()

	// Запускаем миграции в базе данных
	// TODO: Подумать, как откатывать миграции при ошибках
	log.Info("Запускаем миграции")
	postgreSQLMigrator, err := migrator.NewMigrator(
		migrator.MigratorConfig{
			Migrations:      nil,
			Conn:            pgsql.DB.DB,
			EmbedMigrations: pgsqlMigrations.EmbedMigrationsPgsql,
			Dir:             "pgsql",
			Dialect:         goose.DialectPostgres,
		},
	)
	if err != nil {
		return err
	}
	if err = postgreSQLMigrator.Up(ctx); err != nil {
		return err
	}

	// Регистрируем репозитории
	exporterRepository := exporterRepository.NewExporterRepository(pgsql)

	// Регистрируем сервисы
	exporterService := exporterService.NewExporterService(exporterRepository)

	// Регистрируем эндпоинты
	exporterEndpoint.MountExchangeEndpoint(mosquitto, exporterService)

	r := router.NewRouter()
	r.Handle("/metrics", promhttp.Handler())

	server, err := httpServer.GetDefaultServer(conf.Port.HTTP, r)
	if err != nil {
		return err
	}

	// Создаем wait группу
	eg, ctx := errgroup.WithContext(ctx)

	// Запускаем HTTP-сервер
	eg.Go(func() error {

		log.Info(fmt.Sprintf("Server is listening: %s", conf.Port.HTTP))

		return server.Serve()
	})

	// Запускаем горутину, ожидающую завершение контекста
	eg.Go(func() error {

		// Если контекст завершился, значит процесс убили
		<-ctx.Done()

		// Плавно завершаем работу сервера
		server.Shutdown(ctx)

		return nil
	})

	// Ждем завершения контекста или ошибок в горутинах
	return eg.Wait()
}

func initSingletons(conf config.Config) error {

	// Конфигурируем decimal, чтобы в JSON не было кавычек
	decimal.MarshalJSONWithoutQuotes = true

	return nil
}
