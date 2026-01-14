package config

import (
	"pkg/env"
	"pkg/log"
	"pkg/mosquitto"
	"pkg/trace"

	"pkg/database/pgsql"
)

// Config - общая структура конфига
type Config struct {

	// Адрес для http-сервера
	Port struct {
		HTTP string `env:"LISTEN_HTTP"`
	}

	// Данные базы данных
	Pgsql pgsql.PgsqlConfigEnv

	// Подключение к mosquitto
	Mosquitto mosquitto.MosquittoConfigEnv

	Tracer trace.TracerConfig

	ServiceName string `env:"SERVICE_NAME"`
	Environment string `env:"ENVIRONMENT"`

	Logger log.LoggerSettingsEnv

	StackTraceEnabled bool `env:"STACK_TRACE_ENABLED"`
}

func Load() Config {
	return env.Load[Config]()
}
