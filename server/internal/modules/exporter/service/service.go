package service

import (
	"context"
	"pkg/cache"

	"go.opentelemetry.io/otel"

	exporterModel "server/internal/modules/exporter/model"
	exporterRepository "server/internal/modules/exporter/repository"
)

var tracer = otel.Tracer("/server/internal/modules/exporter/service")

var _ ExporterRepository = new(exporterRepository.ExporterRepository)

type ExporterRepository interface {
	CreateMeteoData(context.Context, []exporterModel.CreateMeteoDataReq) error
}

type ArrgCacheValue struct {
	count int
	sum   float64
}

type ExporterService struct {
	exporterRepository ExporterRepository
	arrgCache          *cache.ItemCache[string, ArrgCacheValue]
}

func NewExporterService(
	exporterRepository ExporterRepository,
) *ExporterService {
	s := &ExporterService{
		exporterRepository: exporterRepository,
		arrgCache:          cache.NewItemCache[string, ArrgCacheValue](),
	}

	go s.flush()

	return s
}
