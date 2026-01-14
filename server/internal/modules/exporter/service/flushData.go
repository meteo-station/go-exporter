package service

import (
	"context"
	"pkg/log"
	"server/internal/modules/exporter/model"
	"server/internal/utils/errors"
	"strings"
	"time"
)

func (s *ExporterService) flush() {
	for {
		time.Sleep(time.Minute)
		if err := s._flush(); err != nil {
			log.Error(err)
		}
	}
}

func (s *ExporterService) _flush() error {

	// Получаем все значения из кэша
	aggrTable := s.arrgCache.PopAll()
	if len(aggrTable) == 0 {
		return errors.NotFound.New("items not found")
	}

	var reqs []model.CreateMeteoDataReq

	// Проходим по каждому значению
	for key, item := range aggrTable {

		parts := strings.Split(key, "-")
		if len(parts) != 2 {
			return errors.InternalServer.New("key must contains 2 parts")
		}

		deviceID := parts[0]
		metricType := parts[1]

		avgMetricValue := item.sum / float64(item.count)

		// Парсим ключ
		reqs = append(reqs, model.CreateMeteoDataReq{
			MetricType: metricType,
			DeviceID:   deviceID,
			Value:      avgMetricValue,
		})
	}

	return s.exporterRepository.CreateMeteoData(context.Background(), reqs)
}
