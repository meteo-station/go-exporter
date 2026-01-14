package model

type CreateMeteoDataReq struct {
	MetricType string
	DeviceID   string
	Value      float64
}
