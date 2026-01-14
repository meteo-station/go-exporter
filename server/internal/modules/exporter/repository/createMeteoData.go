package repository

import (
	"context"
	"server/internal/ddl/meteoDataDDL"
	"server/internal/modules/exporter/model"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// CreateMeteoData добавляет метео данные в бд
func (r *ExporterRepository) CreateMeteoData(ctx context.Context, reqs []model.CreateMeteoDataReq) error {
	ctx, span := tracer.Start(ctx, "сreateMeteoData")
	defer span.End()

	now := time.Now()

	q := sq.
		Insert(meteoDataDDL.Table).
		Columns(
			meteoDataDDL.ColumnID,
			meteoDataDDL.ColumnMetricType,
			meteoDataDDL.ColumnDatetime,
			meteoDataDDL.ColumnDeviceID,
			meteoDataDDL.ColumnValue,
		)

	for _, req := range reqs {
		q = q.Values(
			uuid.New(),
			req.MetricType,
			now,
			req.DeviceID,
			req.Value,
		)
	}

	// Добавляем метрику
	return r.db.Exec(ctx, q)
}
