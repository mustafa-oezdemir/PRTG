package prtg

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (a *Api) GetStatusList() (*PrtgStatusListResponse, error) {
	body, err := a.baseExecuteRequest("status.json", nil)
	if err != nil {
		return nil, err
	}

	var response PrtgStatusListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &response, nil
}

func (a *Api) GetGroups() (*PrtgGroupListResponse, error) {
	params := map[string]string{
		"content": "groups",
		"columns": "active,channel,datetime,device,group,message,objid,priority,sensor,status,tags",
		"count":   "50000",
		"output":  "json",
	}

	body, err := a.baseExecuteRequest("table.json", params)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("empty response from PRTG API")
	}

	log.DefaultLogger.Debug("Raw PRTG response",
		"endpoint", "groups",
		"responseSize", len(body),
		"response", string(body),
	)

	var response PrtgGroupListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	if response.Groups == nil {
		return nil, fmt.Errorf("invalid response structure: groups array is nil")
	}

	return &response, nil
}

func (a *Api) GetDevices(group string) (*PrtgDevicesListResponse, error) {
	if group == "" {
		return nil, fmt.Errorf("group parameter is required")
	}

	params := map[string]string{
		"content":      "devices",
		"columns":      "active,channel,datetime,device,group,message,objid,priority,sensor,status,tags",
		"count":        "50000",
		"filter_group": group,
	}

	body, err := a.baseExecuteRequest("table.json", params)
	if err != nil {
		return nil, err
	}

	var response PrtgDevicesListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

func (a *Api) GetSensors(device string) (*PrtgSensorsListResponse, error) {
	if device == "" {
		return nil, fmt.Errorf("device parameter is required")
	}

	params := map[string]string{
		"content":       "sensors",
		"columns":       "active,channel,datetime,device,group,message,objid,priority,sensor,status,tags",
		"count":         "50000",
		"filter_device": device,
	}

	body, err := a.baseExecuteRequest("table.json", params)
	if err != nil {
		return nil, err
	}

	var response PrtgSensorsListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

func (a *Api) GetChannels(objid string) (*PrtgChannelValueStruct, error) {
	params := map[string]string{
		"content":    "values",
		"id":         objid,
		"columns":    "value_,datetime",
		"usecaption": "true",
	}

	body, err := a.baseExecuteRequest("historicdata.json", params)
	if err != nil {
		return nil, err
	}

	var response PrtgChannelValueStruct
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
