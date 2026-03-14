package repository

import (
	"context"
	"server/internal/ddl/homebridgeAccessoryMappingDDL"
	"server/internal/modules/exporter/model"

	sq "github.com/Masterminds/squirrel"
)

func (r *ExporterRepository) GetHomebridgeAccessoryMappings(ctx context.Context) ([]model.HomebridgeAccessoryMapping, error) {
	ctx, span := tracer.Start(ctx, "getHomebridgeAccessoryMappings")
	defer span.End()

	q := sq.
		Select(
			homebridgeAccessoryMappingDDL.ColumnID,
			homebridgeAccessoryMappingDDL.ColumnDeviceID,
			homebridgeAccessoryMappingDDL.ColumnMetricType,
			homebridgeAccessoryMappingDDL.ColumnAccessoryUniqueID,
			homebridgeAccessoryMappingDDL.ColumnValueMapper,
		).
		From(homebridgeAccessoryMappingDDL.Table)

	var mappings []model.HomebridgeAccessoryMapping
	if err := r.db.Select(ctx, &mappings, q); err != nil {
		return nil, err
	}

	return mappings, nil
}
