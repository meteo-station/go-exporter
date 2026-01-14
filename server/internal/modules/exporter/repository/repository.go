package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/exporter/repository")

type ExporterRepository struct {
	db *sql.DB
}

func NewExporterRepository(db *sql.DB) *ExporterRepository {
	return &ExporterRepository{
		db: db,
	}
}
