package homebridgeAccessoryMappingDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaMeteoStation + "." + "homebridge_accessory_mapping"
	TableWithAlias = Table + " " + alias
	alias          = "ham"
)

const (
	ColumnID                = "id"
	ColumnDeviceID          = "device_id"
	ColumnMetricType        = "metric_type"
	ColumnAccessoryUniqueID = "accessory_unique_id"
	ColumnValueMapper       = "value_mapper"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
