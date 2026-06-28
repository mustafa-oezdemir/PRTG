package stream

import (
	"testing"
	"time"
)

func TestAppendBufferDataTrimsOldestValues(t *testing.T) {
	buffer := &dataBuffer{
		times:  []time.Time{time.Unix(1, 0), time.Unix(2, 0)},
		values: []float64{10, 20},
		size:   3,
	}

	appendBufferData(
		buffer,
		[]time.Time{time.Unix(3, 0), time.Unix(4, 0)},
		[]float64{30, 40},
		buffer.size,
	)

	assertBuffer(t, buffer, []int64{2, 3, 4}, []float64{20, 30, 40})
}

func TestAppendBufferDataKeepsNewestValuesWhenIncomingExceedsMaxSize(t *testing.T) {
	buffer := &dataBuffer{
		times:  []time.Time{time.Unix(1, 0)},
		values: []float64{10},
		size:   2,
	}

	appendBufferData(
		buffer,
		[]time.Time{time.Unix(2, 0), time.Unix(3, 0), time.Unix(4, 0)},
		[]float64{20, 30, 40},
		buffer.size,
	)

	assertBuffer(t, buffer, []int64{3, 4}, []float64{30, 40})
}

func TestUpdateChannelBufferAppendModeUpdatesLastValue(t *testing.T) {
	stream := &activeStream{
		updateMode: "append",
		bufferSize: 3,
	}
	state := &channelState{
		buffer: &dataBuffer{
			times:  []time.Time{time.Unix(1, 0)},
			values: []float64{10},
			size:   3,
		},
	}

	updateChannelBuffer(
		stream,
		state,
		[]time.Time{time.Unix(2, 0), time.Unix(3, 0)},
		[]float64{20, 30},
	)

	if state.lastValue != 30 {
		t.Fatalf("expected last value 30, got %f", state.lastValue)
	}
	assertBuffer(t, state.buffer, []int64{1, 2, 3}, []float64{10, 20, 30})
}

func TestUpdateChannelBufferFullModeReplacesAndTrims(t *testing.T) {
	stream := &activeStream{
		updateMode: "full",
		bufferSize: 2,
	}
	state := &channelState{
		buffer: &dataBuffer{
			size: 2,
		},
	}

	updateChannelBuffer(
		stream,
		state,
		[]time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0)},
		[]float64{10, 20, 30},
	)

	if state.lastValue != 30 {
		t.Fatalf("expected last value 30, got %f", state.lastValue)
	}
	assertBuffer(t, state.buffer, []int64{2, 3}, []float64{20, 30})
}

func assertBuffer(t *testing.T, buffer *dataBuffer, wantUnix []int64, wantValues []float64) {
	t.Helper()

	if len(buffer.times) != len(wantUnix) {
		t.Fatalf("expected %d times, got %d", len(wantUnix), len(buffer.times))
	}
	if len(buffer.values) != len(wantValues) {
		t.Fatalf("expected %d values, got %d", len(wantValues), len(buffer.values))
	}

	for i, want := range wantUnix {
		if got := buffer.times[i].Unix(); got != want {
			t.Fatalf("time[%d]: expected unix %d, got %d", i, want, got)
		}
	}
	for i, want := range wantValues {
		if got := buffer.values[i]; got != want {
			t.Fatalf("value[%d]: expected %f, got %f", i, want, got)
		}
	}
}
