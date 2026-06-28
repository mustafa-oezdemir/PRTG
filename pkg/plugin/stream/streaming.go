package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/schema"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

const (
	DefaultBufferSize = 500

	MaxStreamsPerPanel = 5

	MinUpdateInterval = 1000 * time.Millisecond
	MaxUpdateInterval = 60000 * time.Millisecond
	DefaultInterval   = 5000 * time.Millisecond

	DefaultCacheTime = 5 * time.Second

	InitialBufferCapacity = 32
)

func (s *Service) SubscribeStream(ctx context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	s.logger.Debug("Subscribe to stream", "path", req.Path)

	if !strings.HasPrefix(req.Path, "prtg-stream/") {
		return &backend.SubscribeStreamResponse{Status: backend.SubscribeStreamStatusNotFound}, nil
	}

	var query schema.QueryModel
	if err := json.Unmarshal(req.Data, &query); err != nil {
		s.logger.Error("Invalid subscription data", "error", err)
		return &backend.SubscribeStreamResponse{Status: backend.SubscribeStreamStatusPermissionDenied}, nil
	}

	if query.SensorId == "" || (len(query.ChannelArray) == 0 && query.Channel == "") {
		s.logger.Error("Missing required fields", "sensorId", query.SensorId)
		return &backend.SubscribeStreamResponse{Status: backend.SubscribeStreamStatusPermissionDenied}, nil
	}

	panelId := fmt.Sprintf("%v", query.PanelID)
	if panelStreams := s.getStreamsByPanel(panelId); len(panelStreams) >= MaxStreamsPerPanel {
		s.logger.Warn("Maximum streams reached for panel", "panelId", panelId)
		return &backend.SubscribeStreamResponse{Status: backend.SubscribeStreamStatusPermissionDenied}, nil
	}

	if timeRangeInfo, err := extractTimeRangeInfo(req.Data); err == nil {
		s.logger.Debug("Stream time range",
			"windowSize", time.Duration(timeRangeInfo.To-timeRangeInfo.From)*time.Millisecond)
	}

	return &backend.SubscribeStreamResponse{Status: backend.SubscribeStreamStatusOK}, nil
}

func extractTimeRangeInfo(data []byte) (struct{ From, To int64 }, error) {
	var result struct {
		TimeRange struct {
			From int64 `json:"from"`
			To   int64 `json:"to"`
		} `json:"timeRange"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return struct{ From, To int64 }{0, 0}, err
	}

	return struct{ From, To int64 }{
		From: result.TimeRange.From,
		To:   result.TimeRange.To,
	}, nil
}

func (s *Service) PublishStream(ctx context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return &backend.PublishStreamResponse{Status: backend.PublishStreamStatusPermissionDenied}, nil
}

func (s *Service) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	var query schema.QueryModel
	if err := json.Unmarshal(req.Data, &query); err != nil {
		return fmt.Errorf("failed to parse stream data: %w", err)
	}

	if query.SensorId == "" {
		return fmt.Errorf("missing required field: sensorId")
	}

	channels := getChannels(query)
	if len(channels) == 0 {
		return fmt.Errorf("missing required field: channel or channelArray")
	}

	interval := getBoundedInterval(query.StreamInterval)
	streamID := generateStreamID(query, channels)
	timeRangeFrom, timeRangeTo := getTimeRange(query)
	cacheDuration := getCacheDuration(query)
	bufferSize := getBufferSize(query)

	stream := createStream(query, channels, streamID, timeRangeFrom, timeRangeTo,
		interval, cacheDuration, bufferSize)

	s.logger.Info("Stream starting",
		"streamID", streamID,
		"intervalMs", interval.Milliseconds())

	s.registerNewStream(stream, streamID)

	return s.runStreamLoop(ctx, stream, query, sender, timeRangeFrom, timeRangeTo)
}

func getChannels(query schema.QueryModel) []string {
	if len(query.ChannelArray) > 0 {
		return query.ChannelArray
	}
	if query.Channel != "" {
		return []string{query.Channel}
	}
	return nil
}

func getBoundedInterval(requestedInterval int64) time.Duration {
	if requestedInterval <= 0 {
		return DefaultInterval
	}

	interval := time.Duration(requestedInterval) * time.Millisecond
	if interval < MinUpdateInterval {
		return MinUpdateInterval
	}
	if interval > MaxUpdateInterval {
		return MaxUpdateInterval
	}
	return interval
}

func generateStreamID(query schema.QueryModel, channels []string) string {
	channelKey := strings.Join(channels, "_")
	return fmt.Sprintf("%v_%s_%s_%s",
		query.PanelID,
		query.RefID,
		query.SensorId,
		channelKey)
}

func getTimeRange(query schema.QueryModel) (time.Time, time.Time) {
	now := time.Now()
	if query.From > 0 && query.To > 0 {
		return time.Unix(0, query.From*int64(time.Millisecond)),
			time.Unix(0, query.To*int64(time.Millisecond))
	}
	return now.Add(-30 * time.Minute), now
}

func getCacheDuration(query schema.QueryModel) time.Duration {
	if query.CacheTime > 0 {
		return time.Duration(query.CacheTime) * time.Millisecond
	}

	if query.StreamInterval > 0 {
		cacheDuration := time.Duration(query.StreamInterval/2) * time.Millisecond
		if cacheDuration < time.Second {
			return time.Second
		}
		return cacheDuration
	}
	return DefaultCacheTime
}

func getBufferSize(query schema.QueryModel) int64 {
	if query.BufferSize > 0 {
		return query.BufferSize
	}
	return DefaultBufferSize
}

func createStream(query schema.QueryModel, channels []string, streamID string,
	fromTime, toTime time.Time, interval, cacheDuration time.Duration, bufferSize int64) *activeStream {

	stream := &activeStream{
		sensorId:          query.SensorId,
		channelArray:      channels,
		interval:          interval,
		group:             query.Group,
		device:            query.Device,
		sensor:            query.Sensor,
		includeGroupName:  query.IncludeGroupName,
		includeDeviceName: query.IncludeDeviceName,
		includeSensorName: query.IncludeSensorName,
		fromTime:          fromTime,
		toTime:            toTime,
		lastUpdate:        time.Now().Add(-cacheDuration),
		cacheTime:         cacheDuration,
		isActive:          true,
		refID:             query.RefID,
		streamID:          streamID,
		panelId:           fmt.Sprintf("%v", query.PanelID),
		queryId:           query.RefID,
		multiChannelKey:   strings.Join(channels, "_"),
		channelStates:     make(map[string]*channelState, len(channels)),
		updateChan:        make(chan struct{}, 1),
		updateMode:        query.UpdateMode,
		bufferSize:        bufferSize,
		status: &streamStatus{
			active:    true,
			updating:  false,
			lastError: nil,
		},
		lastDataTimestamp: time.Now().UnixMilli(),
	}

	if stream.updateMode == "" {
		stream.updateMode = "append"
	}

	for _, channelName := range channels {
		stream.channelStates[channelName] = &channelState{
			lastValue: 0,
			isActive:  true,
			buffer: &dataBuffer{
				times:  make([]time.Time, 0, InitialBufferCapacity),
				values: make([]float64, 0, InitialBufferCapacity),
				size:   bufferSize,
			},
		}
	}

	return stream
}
