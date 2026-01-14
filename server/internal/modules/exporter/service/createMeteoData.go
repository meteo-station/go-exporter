package service

import (
	"context"
	"server/internal/modules/exporter/model"
	"strings"
)

// CreateMeteoData обновляет счет по конкретным полям
func (s *ExporterService) CreateMeteoData(ctx context.Context, reqs []model.CreateMeteoDataReq) error {
	ctx, span := tracer.Start(ctx, "сreateMeteoData")
	defer span.End()

	// Забираем первый deviceID
	var deviceID string
	if len(reqs) != 0 {
		deviceID = reqs[0].DeviceID
	}

	// Для каждой метрики
	for _, req := range reqs {

		// Получаем ключ для аггрегатора значений
		key := strings.Join([]string{deviceID, req.MetricType}, "-")

		// Получаем аггрегационную таблицу для метрики
		arrgTable := s.arrgCache.Get(key)

		// Инкрементируем счетчики
		arrgTable.count++
		arrgTable.sum += req.Value

		// Записываем значение в таблицу
		s.arrgCache.Set(key, arrgTable)
	}

	return nil
}
