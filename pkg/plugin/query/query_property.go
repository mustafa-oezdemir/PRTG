package query

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtgtime"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/schema"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (s *Service) handlePropertyQuery(ctx context.Context, qm schema.QueryModel, property, filterProperty string, baseFrameName string) backend.DataResponse {
	ctx, span := s.tracer.StartSpan(ctx, "handlePropertyQuery")
	backend.Logger.Info("Context", "ctx", ctx)
	defer span.End()

	s.logger.Debug("Processing property query",
		"property", property,
		"filterProperty", filterProperty,
	)

	isRawMode := qm.QueryType == "raw"
	if isRawMode && !strings.HasSuffix(filterProperty, "_raw") {
		filterProperty += "_raw"
		s.logger.Debug("Converting to raw property",
			"original", property,
			"rawProperty", filterProperty,
		)
	}

	var timesRT []time.Time
	var valuesRT []interface{}

	switch property {
	case "group":
		groups, err := s.api.GetGroups()
		if err != nil {
			return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("API request failed: %v", err))
		}
		for _, g := range groups.Groups {
			if g.Group == qm.Group {
				timestamp, _, err := prtgtime.ParseDateTime(g.Datetime)
				if err != nil {
					continue
				}

				var value interface{}
				switch filterProperty {
				case "active", "active_raw":
					value = selectRawOrFormatted(isRawMode, g.ActiveRAW, g.Active)
				case "message", "message_raw":
					value = selectRawOrFormatted(isRawMode, g.MessageRAW, cleanMessageHTML(g.Message))
				case "priority", "priority_raw":
					value = selectRawOrFormatted(isRawMode, g.PriorityRAW, g.Priority)
				case "status", "status_raw":
					value = selectRawOrFormatted(isRawMode, g.StatusRAW, g.Status)
				case "tags", "tags_raw":
					value = selectRawOrFormatted(isRawMode, g.TagsRAW, g.Tags)
				}

				if value != nil {
					timesRT = append(timesRT, timestamp.UTC())
					valuesRT = append(valuesRT, value)
				}
			}
		}
	case "device":
		if qm.Group == "" {
			return backend.ErrDataResponse(backend.StatusBadRequest, "group parameter is required for device query")
		}
		devices, err := s.api.GetDevices(qm.Group)
		if err != nil {
			return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("API request failed: %v", err))
		}
		for _, dev := range devices.Devices {
			if dev.Device == qm.Device {
				timestamp, _, err := prtgtime.ParseDateTime(dev.Datetime)
				if err != nil {
					continue
				}

				var value interface{}
				switch filterProperty {
				case "active":
					value = dev.Active
				case "active_raw":
					value = dev.ActiveRAW
				case "message":
					value = cleanMessageHTML(dev.Message)
				case "message_raw":
					value = dev.MessageRAW
				case "priority":
					value = dev.Priority
				case "priority_raw":
					value = dev.PriorityRAW
				case "status":
					value = dev.Status
				case "status_raw":
					value = dev.StatusRAW
				case "tags":
					value = dev.Tags
				case "tags_raw":
					value = dev.TagsRAW
				}

				if value != nil {
					timesRT = append(timesRT, timestamp)
					valuesRT = append(valuesRT, value)
				}
			}
		}

	case "sensor":
		if qm.Device == "" {
			return backend.ErrDataResponse(backend.StatusBadRequest, "device parameter is required for sensor query")
		}
		sensors, err := s.api.GetSensors(qm.Device)
		if err != nil {
			return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("API request failed: %v", err))
		}

		for _, s := range sensors.Sensors {
			if s.Sensor == qm.Sensor {
				timestamp, _, err := prtgtime.ParseDateTime(s.Datetime)
				if err != nil {
					continue
				}

				var value interface{}
				switch filterProperty {
				case "status", "status_raw":
					if filterProperty == "status_raw" {
						value = float64(s.StatusRAW)
					} else {
						value = s.Status
					}
				case "active", "active_raw":
					if filterProperty == "active_raw" {
						value = float64(s.ActiveRAW)
					} else {
						value = s.Active
					}
				case "priority", "priority_raw":
					if filterProperty == "priority_raw" {
						value = float64(s.PriorityRAW)
					} else {
						value = s.Priority
					}
				case "message", "message_raw":
					if filterProperty == "message_raw" {
						value = s.MessageRAW
					} else {
						value = cleanMessageHTML(s.Message)
					}
				case "tags", "tags_raw":
					if filterProperty == "tags_raw" {
						value = s.TagsRAW
					} else {
						value = s.Tags
					}
				}

				if value != nil {
					timesRT = []time.Time{timestamp}
					valuesRT = []interface{}{value}
					break
				}
			}
		}
	}

	frameName := fmt.Sprintf("%s_%s_%s", baseFrameName, qm.Property, filterProperty)

	displayName := qm.Property
	if qm.IncludeGroupName && qm.Group != "" {
		displayName = fmt.Sprintf("%s - %s", qm.Group, displayName)
	}
	if qm.IncludeDeviceName && qm.Device != "" {
		displayName = fmt.Sprintf("%s - %s", qm.Device, displayName)
	}
	if qm.IncludeSensorName && qm.Sensor != "" {
		displayName = fmt.Sprintf("%s - %s", qm.Sensor, displayName)
	}
	displayName = fmt.Sprintf("%s (%s)", displayName, filterProperty)

	frame := createPropertyFrameWithDisplayName(timesRT, valuesRT, frameName, displayName)

	return backend.DataResponse{
		Frames: []*data.Frame{frame},
	}
}

func createPropertyFrameWithDisplayName(times []time.Time, values []interface{}, frameName, displayName string) *data.Frame {
	if len(times) == 0 || len(values) == 0 {
		return data.NewFrame(frameName + "_empty")
	}

	timeField := data.NewField("Time", nil, times)
	var valueField *data.Field

	switch values[0].(type) {
	case float64, int:
		floatVals := make([]float64, len(values))
		for i, v := range values {
			switch tv := v.(type) {
			case float64:
				floatVals[i] = tv
			case int:
				floatVals[i] = float64(tv)
			}
		}
		valueField = data.NewField("Value", nil, floatVals)
	case string:
		strVals := make([]string, len(values))
		for i, v := range values {
			strVals[i] = v.(string)
		}
		valueField = data.NewField("Value", nil, strVals)
	default:
		strVals := make([]string, len(values))
		for i, v := range values {
			strVals[i] = fmt.Sprintf("%v", v)
		}
		valueField = data.NewField("Value", nil, strVals)
	}

	valueField.Config = &data.FieldConfig{
		DisplayName: displayName,
	}

	return data.NewFrame(frameName, timeField, valueField)
}

func (s *Service) GetPropertyValue(property string, item interface{}) string {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	isRawRequest := strings.HasSuffix(property, "_raw")
	baseProperty := strings.TrimSuffix(property, "_raw")
	fieldName := cases.Title(language.English).String(baseProperty)

	if isRawRequest {
		fieldName += "_raw"
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		alternatives := []string{
			baseProperty,
			baseProperty + "_raw",
			strings.ToLower(fieldName),
			strings.ToUpper(fieldName),
			baseProperty + "_RAW",
		}

		for _, alt := range alternatives {
			if f := v.FieldByName(alt); f.IsValid() {
				field = f
				break
			}
		}
	}

	if !field.IsValid() {
		return "Unknown"
	}

	val := field.Interface()
	switch v := val.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if isRawRequest {
			if v {
				return "1"
			}
			return "0"
		}
		return strconv.FormatBool(v)
	case string:
		if !isRawRequest && baseProperty == "message" {
			return cleanMessageHTML(v)
		}
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func cleanMessageHTML(message string) string {
	message = strings.ReplaceAll(message, `<div class="status">`, "")
	message = strings.ReplaceAll(message, `<div class="moreicon">`, "")
	message = strings.ReplaceAll(message, "</div>", "")
	return strings.TrimSpace(message)
}

func selectRawOrFormatted(isRaw bool, rawValue, formattedValue interface{}) interface{} {
	if isRaw {
		return rawValue
	}
	return formattedValue
}
