package meteoDataDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaMeteoStation + "." + "meteo_data"
	TableWithAlias = Table + " " + alias
	alias          = "md"
)

const (
	ColumnID         = "id"
	ColumnMetricType = "metric_type"
	ColumnDatetime   = "datetime"
	ColumnDeviceID   = "device_id"
	ColumnValue      = "value"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
