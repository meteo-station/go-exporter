package endpoint

import (
	"context"
	"encoding/json"
	"pkg/log"
	"server/internal/modules/exporter/model"
	"server/internal/utils/errors"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (e *endpoint) createMeteoData(_ mqtt.Client, message mqtt.Message) {

	topic := message.Topic()

	// Второе значение в пути топика это идентификатор устройства
	var deviceID string
	parts := strings.Split(topic, "/")
	if len(parts) >= 2 {
		deviceID = parts[1]
	}

	// Ключ - значение
	keyValue := make(map[string]float64)

	// Десериализуем сообщение
	if err := json.Unmarshal(message.Payload(), &keyValue); err != nil {
		log.Error(errors.InternalServer.Wrap(err))
		return
	}

	reqs := make([]model.CreateMeteoDataReq, 0, len(keyValue))

	// Формируем запись на каждую метрику
	for k, v := range keyValue {
		reqs = append(reqs, model.CreateMeteoDataReq{
			MetricType: k,
			DeviceID:   deviceID,
			Value:      v,
		})
	}

	// Вызываем сервис
	if err := e.service.CreateMeteoData(context.Background(), reqs); err != nil {
		log.Error(err)
		return
	}
}
