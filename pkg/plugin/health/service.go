package health

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/models"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtg"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Service struct {
	api         prtg.PRTGAPI
	logger      observability.PrtgLogger
	clearCaches func()
}

func NewService(api prtg.PRTGAPI, logger observability.PrtgLogger, clearCaches func()) *Service {
	return &Service{
		api:         api,
		logger:      logger,
		clearCaches: clearCaches,
	}
}

func (s *Service) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if s.clearCaches != nil {
		s.clearCaches()
	}

	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	s.logger.Debug("Starting health check")

	status, err := s.api.GetStatusList()
	if err != nil {
		s.logger.Error("PRTG health check failed", "error", err)
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("PRTG API error: %s", err.Error()),
		}, nil
	}

	config, err := models.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)
	timezone := ""
	if err == nil {
		timezone = config.Timezone
	}

	details := map[string]interface{}{
		"version":      status.Version,
		"totalSensors": status.TotalSens,
	}
	if timezone != "" {
		details["timezone"] = timezone
	}

	message := fmt.Sprintf("Data source is working. PRTG Version: %s", status.Version)
	if timezone != "" {
		message = fmt.Sprintf("Data source is working. PRTG Version: %s | Timezone: %s", status.Version, timezone)
	}

	detailsJSON, _ := json.Marshal(details)

	return &backend.CheckHealthResult{
		Status:      backend.HealthStatusOk,
		Message:     message,
		JSONDetails: detailsJSON,
	}, nil
}
