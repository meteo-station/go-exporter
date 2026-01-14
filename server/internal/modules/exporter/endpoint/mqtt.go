package endpoint

import (
	"server/internal/modules/exporter/model"
	exporterService "server/internal/modules/exporter/service"

	"golang.org/x/net/context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type endpoint struct {
	service ExporterService
}

var _ ExporterService = new(exporterService.ExporterService)

type ExporterService interface {
	CreateMeteoData(context.Context, []model.CreateMeteoDataReq) error
}

func MountExchangeEndpoint(mqtt mqtt.Client, service ExporterService) {

	e := &endpoint{
		service: service,
	}

	mqtt.Subscribe("esp-meteo-station/+/data", 0, e.createMeteoData)
}
