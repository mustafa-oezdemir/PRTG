package plugin

import (
	"context"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (d *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	ctx, span := d.tracer.StartSpan(ctx, "CallResourceAdapter")
	defer span.End()

	logger := d.logger.WithContext(log.WithContextualAttributes(ctx, []any{
		"endpoint", "callResource",
		"path", req.Path,
		"method", req.Method,
	}))
	logger.Debug("Resource request received")

	if err := d.resource.CallResource(ctx, req, sender); err != nil {
		d.metrics.IncEndpointRequest("callResource", "error")
		observability.RecordError(span, err, "CallResource failed")
		logger.Error("Resource request failed", "error", err)
		return err
	}

	d.metrics.IncEndpointRequest("callResource", "ok")
	logger.Debug("Resource request completed")
	return nil
}
