package query

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtgtime"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/schema"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func (s *Service) handleMetricsQuery(ctx context.Context, qm schema.QueryModel, timeRange backend.TimeRange, baseFrameName string) backend.DataResponse {
	_, span := s.tracer.StartSpan(ctx, "handleMetricsQuery")
	defer span.End()

	queryStart := time.Now()
	s.logger.Debug("Fetching historical data",
		"sensorId", qm.SensorId,
		"timeRange", fmt.Sprintf("%v to %v", timeRange.From, timeRange.To),
		"channels", qm.ChannelArray,
	)

	response := backend.DataResponse{
		Frames: make([]*data.Frame, 0),
	}

	historicalData, err := s.api.GetHistoricalData(qm.SensorId, timeRange.From.UTC(), timeRange.To.UTC())
	if err != nil {
		s.logger.Error("Failed to fetch historical data",
			"error", err,
			"sensorId", qm.SensorId,
		)
		s.metrics.IncError("historical_data_fetch")
		observability.RecordError(span, err, "Failed to fetch historical data")
		return backend.ErrDataResponse(backend.StatusInternal, "failed to fetch data")
	}

	if len(qm.ChannelArray) == 0 && qm.Channel == "" {
		s.logger.Error("No channels specified")
		s.metrics.IncError("missing_channel")
		return backend.ErrDataResponse(backend.StatusBadRequest, "channel selection required")
	}

	channels := qm.ChannelArray
	if len(channels) == 0 && qm.Channel != "" {
		channels = []string{qm.Channel}
	}

	if len(channels) > 1 {
		timesM := make([]time.Time, 0)
		channelData := make(map[string][]float64)

		for _, channelName := range channels {
			channelData[channelName] = make([]float64, 0)
		}

		if historicalData != nil && len(historicalData.HistData) > 0 {
			for _, item := range historicalData.HistData {
				parsedTime, _, err := prtgtime.ParseDateTime(item.Datetime)
				if err != nil {
					continue
				}

				hasData := false
				tempValues := make(map[string]float64)

				for _, channelName := range channels {
					if val, exists := item.Value[channelName]; exists {
						var floatVal float64
						switch v := val.(type) {
						case float64:
							floatVal = v
						case string:
							if parsed, err := strconv.ParseFloat(v, 64); err == nil {
								floatVal = parsed
							} else {
								continue
							}
						default:
							continue
						}
						tempValues[channelName] = floatVal
						hasData = true
					}
				}

				if hasData {
					timesM = append(timesM, parsedTime)
					for _, channelName := range channels {
						if val, exists := tempValues[channelName]; exists {
							channelData[channelName] = append(channelData[channelName], val)
						} else {
							channelData[channelName] = append(channelData[channelName], 0)
						}
					}
				}
			}
		}

		fields := []*data.Field{
			data.NewField("Time", nil, timesM),
		}
		for _, channelName := range channels {
			displayName := channelName
			if qm.IncludeGroupName && qm.Group != "" {
				displayName = fmt.Sprintf("%s - %s", qm.Group, displayName)
			}
			if qm.IncludeDeviceName && qm.Device != "" {
				displayName = fmt.Sprintf("%s - %s", qm.Device, displayName)
			}
			if qm.IncludeSensorName && qm.Sensor != "" {
				displayName = fmt.Sprintf("%s - %s", qm.Sensor, displayName)
			}

			field := data.NewField(channelName, nil, channelData[channelName]).SetConfig(&data.FieldConfig{
				DisplayName: displayName,
				Custom: map[string]interface{}{
					"refId":     baseFrameName,
					"channel":   channelName,
					"queryType": "multi-channel",
				},
			})
			fields = append(fields, field)
		}

		frame := data.NewFrame(fmt.Sprintf("%s_multi", baseFrameName), fields...)
		frame.Meta = &data.FrameMeta{
			Type: data.FrameTypeTimeSeriesMulti,
			Custom: map[string]interface{}{
				"from":      timeRange.From.UnixMilli(),
				"to":        timeRange.To.UnixMilli(),
				"channels":  channels,
				"stable":    true,
				"duration":  timeRange.To.Sub(timeRange.From).String(),
				"timezone":  "UTC",
				"queryType": "multi-channel",
				"refId":     baseFrameName,
			},
		}

		response.Frames = append(response.Frames, frame)
	} else {
		channelName := channels[0]
		timesM := make([]time.Time, 0)
		valuesM := make([]float64, 0)

		if historicalData != nil && len(historicalData.HistData) > 0 {
			for _, item := range historicalData.HistData {
				parsedTime, _, err := prtgtime.ParseDateTime(item.Datetime)
				if err != nil {
					continue
				}

				if val, exists := item.Value[channelName]; exists {
					var floatVal float64
					switch v := val.(type) {
					case float64:
						floatVal = v
					case string:
						if parsed, err := strconv.ParseFloat(v, 64); err == nil {
							floatVal = parsed
						} else {
							continue
						}
					default:
						continue
					}

					timesM = append(timesM, parsedTime)
					valuesM = append(valuesM, floatVal)
				}
			}
		}

		displayName := channelName
		if qm.IncludeGroupName && qm.Group != "" {
			displayName = fmt.Sprintf("%s - %s", qm.Group, displayName)
		}
		if qm.IncludeDeviceName && qm.Device != "" {
			displayName = fmt.Sprintf("%s - %s", qm.Device, displayName)
		}
		if qm.IncludeSensorName && qm.Sensor != "" {
			displayName = fmt.Sprintf("%s - %s", qm.Sensor, displayName)
		}

		frame := data.NewFrame(fmt.Sprintf("%s_single", baseFrameName),
			data.NewField("Time", nil, timesM),
			data.NewField("Value", nil, valuesM).SetConfig(&data.FieldConfig{
				DisplayName: displayName,
				Custom: map[string]interface{}{
					"refId":     baseFrameName,
					"channel":   channelName,
					"queryType": "single-channel",
				},
			}),
		)

		frame.Meta = &data.FrameMeta{
			Type: data.FrameTypeTimeSeriesMulti,
			Custom: map[string]interface{}{
				"from":      timeRange.From.UnixMilli(),
				"to":        timeRange.To.UnixMilli(),
				"channel":   channelName,
				"stable":    true,
				"duration":  timeRange.To.Sub(timeRange.From).String(),
				"timezone":  "UTC",
				"queryType": "single-channel",
				"refId":     baseFrameName,
			},
		}

		response.Frames = append(response.Frames, frame)
	}

	if len(response.Frames) == 0 {
		response.Frames = append(response.Frames, data.NewFrame(fmt.Sprintf("%s_empty", baseFrameName)))
	}

	duration := time.Since(queryStart)
	s.metrics.ObserveAPILatency("historical_data", duration.Seconds())

	return response
}
