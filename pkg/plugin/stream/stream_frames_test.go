package stream

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func TestExtractLatestStreamingPointUsesCurrentTime(t *testing.T) {
	oldTime := time.Now().Add(-time.Hour)
	frame := data.NewFrame("stream_A_single",
		data.NewField("Time", nil, []time.Time{oldTime}),
		data.NewField("Value", nil, []float64{42}),
	)

	before := time.Now()
	times, values := extractLatestStreamingPoint(frame, 1)
	after := time.Now()

	if len(times) != 1 || len(values) != 1 {
		t.Fatalf("expected one streaming point, got %d times and %d values", len(times), len(values))
	}
	if values[0] != 42 {
		t.Fatalf("expected latest value 42, got %f", values[0])
	}
	if times[0].Before(before) || times[0].After(after) {
		t.Fatalf("expected streaming timestamp between %v and %v, got %v", before, after, times[0])
	}
}

func TestExtractFieldChannelNameUsesMultiChannelFieldName(t *testing.T) {
	frame := data.NewFrame("stream_A_multi",
		data.NewField("Time", nil, []time.Time{time.Now()}),
		data.NewField("CPU-Last", nil, []float64{12}),
		data.NewField("Memory", nil, []float64{34}),
	)

	if got := extractFieldChannelName(frame, 1); got != "CPU-Last" {
		t.Fatalf("expected CPU-Last, got %q", got)
	}
	if got := extractFieldChannelName(frame, 2); got != "Memory" {
		t.Fatalf("expected Memory, got %q", got)
	}
}
