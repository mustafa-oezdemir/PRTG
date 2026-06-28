package stream

import (
	"sync"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type streamManager struct {
	streams          map[string]*activeStream
	mu               sync.RWMutex
	defaultCacheTime time.Duration
	activeStreams    map[string]map[string]*activeStream
}

type streamStatus struct {
	active    bool
	updating  bool
	lastError error
}

type channelState struct {
	lastValue float64
	isActive  bool
	buffer    *dataBuffer
}

type dataBuffer struct {
	times  []time.Time
	values []float64
	size   int64
}

type activeStream struct {
	sensorId          string
	channelArray      []string
	interval          time.Duration
	lastUpdate        time.Time
	group             string
	device            string
	sensor            string
	includeGroupName  bool
	includeDeviceName bool
	includeSensorName bool
	fromTime          time.Time
	toTime            time.Time
	cacheTime         time.Duration
	timeRange         *backend.TimeRange
	isActive          bool
	updateChan        chan struct{}
	status            *streamStatus
	refID             string
	streamID          string
	panelId           string
	queryId           string
	multiChannelKey   string
	channelStates     map[string]*channelState
	updateMode        string
	bufferSize        int64
	errorCount        int
	lastDataTimestamp int64
}
