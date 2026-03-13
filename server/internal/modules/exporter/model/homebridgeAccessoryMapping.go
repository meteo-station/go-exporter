package model

type HomebridgeAccessoryMapping struct {
	ID                string `db:"id"`
	DeviceID          string `db:"device_id"`
	MetricType        string `db:"metric_type"`
	AccessoryUniqueID string `db:"accessory_unique_id"`
}