package stream

import (
	"context"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/schema"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type MetricsQueryHandler interface {
	HandleMetricsQuery(ctx context.Context, qm schema.QueryModel, timeRange backend.TimeRange, baseFrameName string) backend.DataResponse
}

type Service struct {
	logger        observability.PrtgLogger
	metricsQuery  MetricsQueryHandler
	streamManager *streamManager
}

func NewService(logger observability.PrtgLogger, metricsQuery MetricsQueryHandler, defaultCacheTime time.Duration) *Service {
	return &Service{
		logger:       logger,
		metricsQuery: metricsQuery,
		streamManager: &streamManager{
			streams:          make(map[string]*activeStream),
			activeStreams:    make(map[string]map[string]*activeStream),
			defaultCacheTime: defaultCacheTime,
		},
	}
}
