package plugin

import (
	"context"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (d *Datasource) SubscribeStream(ctx context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	ctx, span := d.tracer.StartSpan(ctx, "SubscribeStream")
	defer span.End()

	logger := d.logger.WithContext(log.WithContextualAttributes(ctx, []any{
		"endpoint", "subscribeStream",
		"path", req.Path,
	}))
	logger.Debug("Subscribe stream request received")

	response, err := d.stream.SubscribeStream(ctx, req)
	if err != nil {
		d.metrics.IncEndpointRequest("subscribeStream", "error")
		observability.RecordError(span, err, "SubscribeStream failed")
		logger.Error("Subscribe stream request failed", "error", err)
		return nil, err
	}

	d.metrics.IncEndpointRequest("subscribeStream", "ok")
	logger.Debug("Subscribe stream request completed")
	return response, nil
}

func (d *Datasource) PublishStream(ctx context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	ctx, span := d.tracer.StartSpan(ctx, "PublishStream")
	defer span.End()

	logger := d.logger.WithContext(log.WithContextualAttributes(ctx, []any{
		"endpoint", "publishStream",
		"path", req.Path,
	}))
	logger.Debug("Publish stream request received")

	response, err := d.stream.PublishStream(ctx, req)
	if err != nil {
		d.metrics.IncEndpointRequest("publishStream", "error")
		observability.RecordError(span, err, "PublishStream failed")
		logger.Error("Publish stream request failed", "error", err)
		return nil, err
	}

	d.metrics.IncEndpointRequest("publishStream", "ok")
	logger.Debug("Publish stream request completed")
	return response, nil
}

func (d *Datasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	ctx, span := d.tracer.StartSpan(ctx, "RunStream")
	defer span.End()

	logger := d.logger.WithContext(log.WithContextualAttributes(ctx, []any{
		"endpoint", "runStream",
		"path", req.Path,
	}))
	logger.Debug("Run stream request received")

	if err := d.stream.RunStream(ctx, req, sender); err != nil {
		d.metrics.IncEndpointRequest("runStream", "error")
		observability.RecordError(span, err, "RunStream failed")
		logger.Error("Run stream request failed", "error", err)
		return err
	}

	d.metrics.IncEndpointRequest("runStream", "ok")
	logger.Debug("Run stream request completed")
	return nil
}
