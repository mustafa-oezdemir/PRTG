package prtg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (a *Api) GetHistoricalData(sensorID string, startDate, endDate time.Time) (*PrtgHistoricalDataResponse, error) {
	if sensorID == "" {
		return nil, fmt.Errorf("invalid query: missing sensor ID")
	}

	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		loc = time.Local
	}

	localStartDate := startDate.In(loc).Add(-1 * time.Hour)
	localEndDate := endDate.In(loc).Add(1 * time.Hour)

	const format = "2006-01-02-15-04-05"
	sdate := localStartDate.Format(format)
	edate := localEndDate.Format(format)

	hours := localEndDate.Sub(localStartDate).Hours()

	var avg string
	switch {
	case hours <= 12:
		avg = "0"
	case hours <= 24:
		avg = "120"
	case hours <= 48:
		avg = "300"
	case hours <= 96:
		avg = "600"
	case hours <= 168:
		avg = "900"
	case hours <= 336:
		avg = "1800"
	case hours <= 720:
		avg = "3600"
	case hours <= 1440:
		avg = "7200"
	case hours <= 2880:
		avg = "14400"
	case hours <= 4320:
		avg = "28800"
	case hours <= 10080:
		avg = "43200"
	case hours <= 20160:
		avg = "57600"
	case hours <= 43200:
		avg = "86400"
	default:
		avg = "172800"
	}

	params := map[string]string{
		"id":         sensorID,
		"columns":    "datetime,value_",
		"avg":        avg,
		"sdate":      sdate,
		"edate":      edate,
		"count":      "50000",
		"usecaption": "1",
	}

	log.DefaultLogger.Debug("Requesting historical data",
		"sensorID", sensorID,
		"startDate", sdate,
		"endDate", edate,
		"avg", avg,
	)

	cacheKey := fmt.Sprintf("hist_%s_%s_%s", sensorID, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))

	a.cacheMu.RLock()
	if cached, exists := a.cache[cacheKey]; exists && time.Now().Before(cached.expiry) {
		a.cacheMu.RUnlock()
		var response PrtgHistoricalDataResponse
		if err := json.Unmarshal(cached.data, &response); err == nil {
			return &response, nil
		}
	}
	a.cacheMu.RUnlock()

	body, err := a.baseExecuteRequest("historicdata.json", params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical data: %w", err)
	}

	var response PrtgHistoricalDataResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.HistData) == 0 {
		log.DefaultLogger.Debug("No data found for the given time range",
			"sensorID", sensorID,
			"startDate", startDate,
			"endDate", endDate,
		)
		return &response, nil
	}

	if len(response.HistData) > 0 {
		a.cacheMu.Lock()
		if data, err := json.Marshal(response); err == nil {
			a.cache[cacheKey] = cacheItem{
				data:   data,
				expiry: time.Now().Add(a.cacheTime),
			}
		}
		a.cacheMu.Unlock()
	}

	return &response, nil
}
