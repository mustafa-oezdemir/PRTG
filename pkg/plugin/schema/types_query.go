package schema

import (
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type QueryModel struct {
	QueryType         string   `json:"queryType"`
	SensorId          string   `json:"sensorId"`
	DeviceId          string   `json:"deviceId"`
	GroupId           string   `json:"groupId"`
	Group             string   `json:"group"`
	Device            string   `json:"device"`
	Sensor            string   `json:"sensor"`
	Channel           string   `json:"channel"`
	ChannelArray      []string `json:"channelArray"`
	Property          string   `json:"property"`
	FilterProperty    string   `json:"filterProperty"`
	IncludeGroupName  bool     `json:"includeGroupName"`
	IncludeDeviceName bool     `json:"includeDeviceName"`
	IncludeSensorName bool     `json:"includeSensorName"`
	From              int64    `json:"from"`
	To                int64    `json:"to"`
	ManualMethod      string   `json:"manualMethod"`
	ManualObjectId    string   `json:"manualObjectId"`
	Limit             int64    `json:"limit"`
	Tags              []string `json:"tags"`
	DashboardID       int64    `json:"dashboardId"`
	DashboardUID      string   `json:"dashboardUid"`
	PanelID           int64    `json:"panelId"`
	IsStreaming       bool     `json:"isStreaming"`
	StreamInterval    int64    `json:"streamInterval"`
	CacheTime         int64    `json:"cacheTime"`
	BufferSize        int64    `json:"bufferSize"`
	UpdateMode        string   `json:"updateMode"`
	RefID             string   `json:"refId"`
}

type QueryCacheKey struct {
	RefID      string
	QueryType  string
	SensorID   string
	Channel    string
	TimeRange  string
	Property   string
	Parameters string
}

func (k QueryCacheKey) String() string {
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s",
		k.RefID,
		k.QueryType,
		k.SensorID,
		k.Channel,
		k.TimeRange,
		k.Property,
		k.Parameters,
	)
}

type QueryCacheEntry struct {
	Response   backend.DataResponse
	ValidUntil time.Time
	Updating   bool
}
