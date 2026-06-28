package plugin

import (
	"context"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	ctx, span := d.tracer.StartSpan(ctx, "CheckHealth")
	defer span.End()

	logger := d.logger.WithContext(log.WithContextualAttributes(ctx, []any{
		"endpoint", "checkHealth",
	}))
	logger.Debug("Health check request received")

	result, err := d.health.CheckHealth(ctx, req)
	if err != nil {
		d.metrics.IncEndpointRequest("checkHealth", "error")
		observability.RecordError(span, err, "CheckHealth failed")
		logger.Error("Health check request failed", "error", err)
		return nil, err
	}

	status := "ok"
	if result != nil && result.Status != backend.HealthStatusOk {
		status = "error"
	}
	d.metrics.IncEndpointRequest("checkHealth", status)
	logger.Debug("Health check request completed", "status", status)
	return result, nil
}
