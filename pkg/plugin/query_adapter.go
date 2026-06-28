package plugin

import (
	"context"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx, span := d.tracer.StartSpan(ctx, "QueryData")
	defer span.End()

	logger := d.logger.WithContext(log.WithContextualAttributes(ctx, []any{
		"endpoint", "queryData",
		"queryCount", len(req.Queries),
	}))
	logger.Debug("Query data request received")

	for _, q := range req.Queries {
		d.metrics.IncQuery(q.QueryType)
	}

	response, err := d.query.QueryData(ctx, req)
	if err != nil {
		d.metrics.IncEndpointRequest("queryData", "error")
		observability.RecordError(span, err, "QueryData failed")
		logger.Error("Query data request failed", "error", err)
		return nil, err
	}

	d.metrics.IncEndpointRequest("queryData", "ok")
	logger.Debug("Query data request completed")
	return response, nil
}
