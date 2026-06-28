package observability

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	logger PrtgLogger
}

func NewTracer(logger PrtgLogger) *Tracer {
	return &Tracer{
		logger: logger,
	}
}

func (t *Tracer) StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	t.logger.Debug("Starting span", "name", name)
	return tracing.DefaultTracer().Start(ctx, fmt.Sprintf("PRTG.%s", name))
}

func (t *Tracer) AddAttribute(span trace.Span, key string, value interface{}) {
	span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", value)))
}

// recordError adds error details to a span
func RecordError(span trace.Span, err error, message string) {
	span.RecordError(err)
	span.SetAttributes(
		attribute.String("error.message", err.Error()),
		attribute.String("error.type", fmt.Sprintf("%T", err)),
	)
	span.SetStatus(codes.Error, message)
}
