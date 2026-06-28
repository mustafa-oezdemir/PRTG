package query

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/schema"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/concurrent"
)

const MaxConcurrentQueries = 10

func (s *Service) handleSingleQueryData(ctx context.Context, q concurrent.Query) backend.DataResponse {
	ctx, span := s.tracer.StartSpan(ctx, "handleSingleQueryData")
	defer span.End()

	start := time.Now()
	s.logger.Debug("Processing concurrent query",
		"refId", q.DataQuery.RefID,
	)

	res := s.query(ctx, q.PluginContext, q.DataQuery)

	duration := time.Since(start)
	s.metrics.ObserveQueryDuration("single_query", duration.Seconds())
	s.logger.Debug("Concurrent query processed",
		"refId", q.DataQuery.RefID,
		"duration", duration,
		"status", res.Status,
		"hasError", res.Error != nil,
	)

	return res
}

func generateCacheKey(req *backend.QueryDataRequest) string {
	var keyBuilder strings.Builder
	for _, q := range req.Queries {
		keyBuilder.WriteString(fmt.Sprintf("%s:%d:%d:%s;",
			q.RefID,
			q.TimeRange.From.UnixNano(),
			q.TimeRange.To.UnixNano(),
			string(q.JSON)))
	}
	return keyBuilder.String()
}

func (s *Service) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	if len(req.Queries) == 0 {
		return &backend.QueryDataResponse{
			Responses: make(map[string]backend.DataResponse),
		}, nil
	}

	if len(req.Queries) > MaxConcurrentQueries {
		return &backend.QueryDataResponse{
			Responses: map[string]backend.DataResponse{
				req.Queries[0].RefID: {
					Error:  fmt.Errorf("query limit exceeded: %d/%d", len(req.Queries), MaxConcurrentQueries),
					Status: backend.StatusTooManyRequests,
				},
			},
		}, nil
	}

	cacheKey := generateCacheKey(req)
	s.cacheMutex.RLock()
	if cached, exists := s.queryCache[cacheKey]; exists && time.Now().Before(cached.ValidUntil) {
		s.cacheMutex.RUnlock()
		response := backend.NewQueryDataResponse()
		response.Responses[req.Queries[0].RefID] = cached.Response
		return response, nil
	}
	s.cacheMutex.RUnlock()

	response, err := s.mux.QueryData(ctx, req)
	if err != nil {
		return nil, err
	}

	s.cacheMutex.Lock()
	s.queryCache[cacheKey] = &schema.QueryCacheEntry{
		Response:   response.Responses[req.Queries[0].RefID],
		ValidUntil: time.Now().Add(s.cacheTime),
		Updating:   false,
	}
	s.cacheMutex.Unlock()

	return response, nil
}

func (s *Service) handleMetricsQueryType(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx, "handleMetricsQueryType")
	defer span.End()

	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		response.Responses[q.RefID] = s.handleSingleQueryData(ctx, concurrent.Query{
			DataQuery:     q,
			PluginContext: req.PluginContext,
		})
	}

	return response, nil
}

func (s *Service) handleManualQueryType(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	_, span := s.tracer.StartSpan(ctx, "handleManualQueryType")
	defer span.End()

	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		var qm schema.QueryModel
		if err := json.Unmarshal(q.JSON, &qm); err != nil {
			response.Responses[q.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, "failed to parse query")
			continue
		}

		response.Responses[q.RefID] = s.handleManualQuery(qm, q.TimeRange, fmt.Sprintf("manual_%s", q.RefID))
	}

	return response, nil
}

func (s *Service) handlePropertyQueryType(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx, "handlePropertyQueryType")
	defer span.End()

	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		var qm schema.QueryModel
		if err := json.Unmarshal(q.JSON, &qm); err != nil {
			response.Responses[q.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, "failed to parse query")
			continue
		}

		response.Responses[q.RefID] = s.handlePropertyQuery(
			ctx,
			qm,
			qm.Property,
			qm.FilterProperty,
			fmt.Sprintf("property_%s", q.RefID),
		)
	}

	return response, nil
}

func (s *Service) handleQueryFallback(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	s.logger.Warn("Query type not supported", "queries", len(req.Queries))
	return backend.NewQueryDataResponse(), nil
}
