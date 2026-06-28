package query

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/schema"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func (s *Service) query(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var qm struct {
		CacheTime int64 `json:"cacheTime"`
	}
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		qm.CacheTime = 6000
	}

	cacheKey := schema.QueryCacheKey{
		RefID:     query.RefID,
		QueryType: query.QueryType,
		TimeRange: fmt.Sprintf("%v-%v", query.TimeRange.From.Unix(), query.TimeRange.To.Unix()),
	}

	s.cacheMutex.RLock()
	if cached, exists := s.queryCache[cacheKey.String()]; exists &&
		time.Now().Before(cached.ValidUntil) {
		s.cacheMutex.RUnlock()
		return cached.Response
	}
	s.cacheMutex.RUnlock()

	response := s.executeQuery(ctx, pCtx, query)

	if response.Error == nil {
		s.cacheMutex.Lock()
		s.queryCache[cacheKey.String()] = &schema.QueryCacheEntry{
			Response:   response,
			ValidUntil: time.Now().Add(time.Duration(qm.CacheTime) * time.Millisecond),
		}
		s.cacheMutex.Unlock()
	}

	return response
}

func (s *Service) executeQuery(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	ctx, span := s.tracer.StartSpan(ctx, "query")
	defer span.End()
	backend.Logger.Info("PluginContext", "pCtx", pCtx)

	start := time.Now()
	s.logger.Info("Starting query execution",
		"queryType", query.QueryType,
		"refID", query.RefID,
		"timeRange", fmt.Sprintf("%v to %v", query.TimeRange.From, query.TimeRange.To),
	)

	var qm schema.QueryModel
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		s.logger.Error("Query parsing failed",
			"error", err,
			"raw_query", string(query.JSON),
		)
		s.metrics.IncError("query_parse_error")
		observability.RecordError(span, err, "Failed to parse query")
		return backend.ErrDataResponse(backend.StatusBadRequest, "failed to parse query")
	}

	cacheKey := schema.QueryCacheKey{
		RefID:      query.RefID,
		QueryType:  query.QueryType,
		SensorID:   qm.SensorId,
		Channel:    strings.Join(qm.ChannelArray, ","),
		TimeRange:  fmt.Sprintf("%v-%v", query.TimeRange.From.Unix(), query.TimeRange.To.Unix()),
		Property:   qm.Property,
		Parameters: fmt.Sprintf("%s_%s_%s", qm.Group, qm.Device, qm.Sensor),
	}

	cacheTime := s.api.GetCacheTime()
	timeRange := query.TimeRange.To.Sub(query.TimeRange.From)
	var cacheDuration time.Duration

	switch {
	case timeRange <= time.Hour:
		cacheDuration = 6 * time.Second
	case timeRange <= 24*time.Hour:
		cacheDuration = 30 * time.Second
	default:
		cacheDuration = cacheTime
	}

	cacheKeyStr := cacheKey.String()

	s.cacheMutex.RLock()
	if entry, exists := s.queryCache[cacheKeyStr]; exists && time.Now().Before(entry.ValidUntil) {
		s.cacheMutex.RUnlock()
		return entry.Response
	}
	s.cacheMutex.RUnlock()

	addQueryAttributes(span, qm)

	defer func() {
		duration := time.Since(start).Seconds()
		s.metrics.ObserveQueryDuration(qm.QueryType, duration)
		s.logger.Info("Query completed",
			"duration", duration,
			"queryType", qm.QueryType,
			"refID", query.RefID,
		)
	}()

	var response backend.DataResponse
	switch qm.QueryType {
	case "metrics":
		if qm.Channel == "" && len(qm.ChannelArray) == 0 {
			s.logger.Error("Channel selection required for metrics query")
			s.metrics.IncError("missing_channel")
			return backend.ErrDataResponse(backend.StatusBadRequest, "channel selection required")
		}
		response = s.handleMetricsQuery(ctx, qm, query.TimeRange, fmt.Sprintf("metrics_%s", query.RefID))

		if response.Error == nil {
			s.cacheMutex.Lock()
			s.queryCache[cacheKey.String()] = &schema.QueryCacheEntry{
				Response:   response,
				ValidUntil: time.Now().Add(25 * time.Second),
				Updating:   false,
			}
			s.cacheMutex.Unlock()
		}

	case "manual":
		s.logger.Debug("Executing manual query",
			"method", qm.ManualMethod,
			"objectId", qm.ManualObjectId,
		)
		response = s.handleManualQuery(qm, query.TimeRange, fmt.Sprintf("manual_%s", query.RefID))

	case "text", "raw":
		response = s.handlePropertyQuery(ctx, qm, qm.Property, qm.FilterProperty, fmt.Sprintf("property_%s", query.RefID))

	default:
		s.logger.Warn("Unknown query type",
			"type", qm.QueryType,
			"refID", query.RefID,
		)
		s.metrics.IncError("unknown_query_type")
		return backend.DataResponse{
			Frames: []*data.Frame{
				data.NewFrame(fmt.Sprintf("unknown_%s", query.RefID)),
			},
		}
	}

	if response.Error == nil {
		s.cacheMutex.Lock()
		s.queryCache[cacheKeyStr] = &schema.QueryCacheEntry{
			Response:   response,
			ValidUntil: time.Now().Add(cacheDuration),
			Updating:   false,
		}
		s.cacheMutex.Unlock()

		s.logger.Debug("Cached response",
			"key", cacheKeyStr,
			"duration", cacheDuration,
		)
	}

	if response.Error != nil {
		s.logger.Error("Query execution failed",
			"error", response.Error,
			"queryType", qm.QueryType,
			"refID", query.RefID,
		)
		s.metrics.IncError("query_execution")
		observability.RecordError(span, response.Error, "Query execution failed")
	}

	return response
}
