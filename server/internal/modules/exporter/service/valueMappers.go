package service

const ValueMapperAirQuality = "air_quality"

// applyValueMapper применяет маппер к значению перед отправкой в Homebridge.
// Если маппер не задан — возвращает исходное значение.
func applyValueMapper(mapper *string, value float64) float64 {
	if mapper == nil {
		return value
	}
	switch *mapper {
	case ValueMapperAirQuality:
		return mapIAQToHomeKit(value)
	default:
		return value
	}
}

// mapIAQToHomeKit переводит IAQ индекс BME688 (0–500) в шкалу HomeKit (1–5):
// 1 = Excellent, 2 = Good, 3 = Fair, 4 = Inferior, 5 = Poor
func mapIAQToHomeKit(iaq float64) float64 {
	switch {
	case iaq <= 50:
		return 1
	case iaq <= 100:
		return 2
	case iaq <= 150:
		return 3
	case iaq <= 200:
		return 4
	default:
		return 5
	}
}