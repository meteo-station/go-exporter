package service

import (
	"context"
	"pkg/cache"

	"go.opentelemetry.io/otel"

	"server/internal/config"
	exporterModel "server/internal/modules/exporter/model"
	exporterRepository "server/internal/modules/exporter/repository"
)

var tracer = otel.Tracer("/server/internal/modules/exporter/service")

var _ ExporterRepository = new(exporterRepository.ExporterRepository)

type ExporterRepository interface {
	CreateMeteoData(context.Context, []exporterModel.CreateMeteoDataReq) error
	GetHomebridgeAccessoryMappings(context.Context) ([]exporterModel.HomebridgeAccessoryMapping, error)
}

type ArrgCacheValue struct {
	count int
	sum   float64
}

type ExporterService struct {
	exporterRepository ExporterRepository
	arrgCache          *cache.ItemCache[string, ArrgCacheValue]
	homebridgeConf     config.HomebridgeConfig
}

func NewExporterService(
	exporterRepository ExporterRepository,
	homebridgeConf config.HomebridgeConfig,
) *ExporterService {
	s := &ExporterService{
		exporterRepository: exporterRepository,
		arrgCache:          cache.NewItemCache[string, ArrgCacheValue](),
		homebridgeConf:     homebridgeConf,
	}

	go s.flush()

	return s
}
