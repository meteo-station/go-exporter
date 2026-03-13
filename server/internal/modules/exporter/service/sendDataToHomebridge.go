package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"server/internal/utils/errors"

	exporterModel "server/internal/modules/exporter/model"
)

func (s *ExporterService) sendDataToHomebridge(ctx context.Context, data []exporterModel.CreateMeteoDataReq) error {
	if s.homebridgeConf.WebhookURL == "" {
		return nil
	}

	mappings, err := s.exporterRepository.GetHomebridgeAccessoryMappings(ctx)
	if err != nil {
		return err
	}

	if len(mappings) == 0 {
		return nil
	}

	// Строим индекс: "deviceID-metricType" -> accessoryID
	mappingIndex := make(map[string]string, len(mappings))
	for _, m := range mappings {
		mappingIndex[m.DeviceID+"-"+m.MetricType] = m.AccessoryUniqueID
	}

	for _, item := range data {
		accessoryID, ok := mappingIndex[item.DeviceID+"-"+item.MetricType]
		if !ok {
			continue
		}
		if err := sendWebhook(ctx, s.homebridgeConf.WebhookURL, accessoryID, item.Value); err != nil {
			return fmt.Errorf("webhook accessory %s: %w", accessoryID, err)
		}
	}

	return nil
}

func sendWebhook(ctx context.Context, webhookURL, accessoryID string, value float64) error {
	u := fmt.Sprintf("%s/?accessoryId=%s&value=%g", webhookURL, url.QueryEscape(accessoryID), value)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.BadGateway.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return errors.BadGateway.New(fmt.Sprintf("homebridge webhook returned status %d", resp.StatusCode))
	}

	return nil
}