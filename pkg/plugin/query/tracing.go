package query

import (
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/schema"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func addQueryAttributes(span trace.Span, query schema.QueryModel) {
	attrs := []attribute.KeyValue{
		attribute.String("query.group", query.Group),
		attribute.String("query.groupId", query.GroupId),
		attribute.String("query.device", query.Device),
		attribute.String("query.deviceId", query.DeviceId),
		attribute.String("query.sensor", query.Sensor),
		attribute.String("query.sensorId", query.SensorId),
	}

	if query.Channel != "" {
		attrs = append(attrs, attribute.String("query.channel", query.Channel))
	}

	if query.Property != "" {
		attrs = append(attrs, attribute.String("query.property", query.Property))
	}

	if query.FilterProperty != "" {
		attrs = append(attrs, attribute.String("query.filterProperty", query.FilterProperty))
	}

	span.SetAttributes(attrs...)
}
