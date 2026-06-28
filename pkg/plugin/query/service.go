package query

import (
	"context"
	"sync"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtg"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/schema"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
)

type Service struct {
	api        prtg.PRTGAPI
	logger     observability.PrtgLogger
	tracer     *observability.Tracer
	metrics    *observability.Metrics
	mux        backend.QueryDataHandler
	queryCache map[string]*schema.QueryCacheEntry
	cacheMutex sync.RWMutex
	cacheTime  time.Duration
}

func NewService(
	api prtg.PRTGAPI,
	logger observability.PrtgLogger,
	tracer *observability.Tracer,
	metrics *observability.Metrics,
	cacheTime time.Duration,
) *Service {
	s := &Service{
		api:        api,
		logger:     logger,
		tracer:     tracer,
		metrics:    metrics,
		queryCache: make(map[string]*schema.QueryCacheEntry),
		cacheMutex: sync.RWMutex{},
		cacheTime:  cacheTime,
	}

	queryTypeMux := datasource.NewQueryTypeMux()
	queryTypeMux.HandleFunc("metrics", s.handleMetricsQueryType)
	queryTypeMux.HandleFunc("manual", s.handleManualQueryType)
	queryTypeMux.HandleFunc("text", s.handlePropertyQueryType)
	queryTypeMux.HandleFunc("raw", s.handlePropertyQueryType)
	queryTypeMux.HandleFunc("", s.handleQueryFallback)

	s.mux = queryTypeMux
	return s
}

func (s *Service) ClearCache() {
	s.cacheMutex.Lock()
	s.queryCache = make(map[string]*schema.QueryCacheEntry)
	s.cacheMutex.Unlock()
}

func (s *Service) HandleMetricsQuery(ctx context.Context, qm schema.QueryModel, timeRange backend.TimeRange, baseFrameName string) backend.DataResponse {
	return s.handleMetricsQuery(ctx, qm, timeRange, baseFrameName)
}
