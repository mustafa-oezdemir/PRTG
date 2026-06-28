package stream

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/schema"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func (s *Service) runStreamLoop(
	ctx context.Context,
	stream *activeStream,
	query schema.QueryModel,
	sender *backend.StreamSender,
	timeRangeFrom, timeRangeTo time.Time,
) error {
	defer s.cleanupStream(stream)

	timeRange := backend.TimeRange{
		From: timeRangeFrom,
		To:   timeRangeTo,
	}

	jitter := time.Duration(rand.Int63n(250)) * time.Millisecond
	ticker := time.NewTicker(stream.interval + jitter)
	defer ticker.Stop()

	if err := s.updateStreamWithMetricsQuery(ctx, stream, sender, query, timeRange); err != nil {
		s.logger.Error("Initial stream update failed", "error", err)
	}

	return s.streamUpdateLoop(ctx, stream, sender, query, timeRange, ticker)
}

func (s *Service) streamUpdateLoop(
	ctx context.Context,
	stream *activeStream,
	sender *backend.StreamSender,
	query schema.QueryModel,
	timeRange backend.TimeRange,
	ticker *time.Ticker,
) error {
	for {
		select {
		case <-ctx.Done():
			stream.isActive = false
			return nil

		case <-ticker.C:
			advanceTimeRange(&timeRange)

			updateCtx, cancel := context.WithTimeout(ctx, getUpdateTimeout(stream.interval))
			err := s.updateStreamWithMetricsQuery(updateCtx, stream, sender, query, timeRange)
			cancel()

			s.handleStreamError(stream, err)

		case <-stream.updateChan:
			if err := s.updateStreamWithMetricsQuery(ctx, stream, sender, query, timeRange); err != nil {
				s.logger.Error("Manual update failed", "error", err)
			}
		}
	}
}

func getUpdateTimeout(interval time.Duration) time.Duration {
	timeout := interval * 2
	if timeout < 10*time.Second {
		return 10 * time.Second
	}
	return timeout
}

func advanceTimeRange(timeRange *backend.TimeRange) {
	now := time.Now()
	windowSize := timeRange.To.Sub(timeRange.From)
	if windowSize <= 0 {
		windowSize = 30 * time.Minute
	}
	timeRange.From = now.Add(-windowSize)
	timeRange.To = now
}

func (s *Service) handleStreamError(stream *activeStream, err error) {
	if err != nil {
		stream.errorCount++
		if stream.errorCount <= 3 || stream.errorCount%10 == 0 {
			s.logger.Error("Stream update failed",
				"error", err,
				"count", stream.errorCount,
				"streamID", stream.streamID)
		}
	} else {
		stream.errorCount = 0
	}
}

func (s *Service) updateStreamWithMetricsQuery(
	ctx context.Context,
	stream *activeStream,
	sender *backend.StreamSender,
	query schema.QueryModel,
	timeRange backend.TimeRange,
) error {
	if time.Since(stream.lastUpdate) < stream.cacheTime {
		return nil
	}

	stream.status.updating = true
	defer func() { stream.status.updating = false }()

	stream.lastUpdate = time.Now()

	baseFrameName := fmt.Sprintf("stream_%s", stream.refID)
	response := s.metricsQuery.HandleMetricsQuery(ctx, query, timeRange, baseFrameName)

	if response.Error != nil {
		stream.status.lastError = response.Error
		return response.Error
	}

	if len(response.Frames) == 0 {
		return nil
	}

	return s.processResponseFrames(stream, sender, response, timeRange)
}

func (s *Service) processResponseFrames(
	stream *activeStream,
	sender *backend.StreamSender,
	response backend.DataResponse,
	timeRange backend.TimeRange,
) error {
	for _, frame := range response.Frames {
		if len(frame.Fields) < 2 {
			continue
		}

		for fieldIndex := 1; fieldIndex < len(frame.Fields); fieldIndex++ {
			channelName := extractFieldChannelName(frame, fieldIndex)
			if channelName == "" {
				continue
			}

			channelState, exists := stream.channelStates[channelName]
			if !exists {
				continue
			}

			times, values := extractLatestStreamingPoint(frame, fieldIndex)
			if len(times) == 0 {
				continue
			}

			updateChannelBuffer(stream, channelState, times, values)

			streamFrame := createStreamingFrame(stream, channelName, channelState, timeRange.From, timeRange.To)

			if err := sender.SendFrame(streamFrame, data.IncludeAll); err != nil {
				s.logger.Error("Failed to send frame", "error", err, "channel", channelName)
				continue
			}
		}
	}

	stream.lastDataTimestamp = time.Now().UnixMilli()
	return nil
}
