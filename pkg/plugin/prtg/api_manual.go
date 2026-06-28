package prtg

import (
	"encoding/json"
	"fmt"
)

func (a *Api) ExecuteManualMethod(method string, objectId string) (*PrtgManualMethodResponse, error) {
	params := map[string]string{}

	if objectId != "" {
		params["id"] = objectId
	}

	body, err := a.baseExecuteRequest(method, params)
	if err != nil {
		return nil, fmt.Errorf("manual API request failed: %w", err)
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var keyValues []KeyValue
	flattenJSON("", rawData, &keyValues)

	return &PrtgManualMethodResponse{
		Manuel:    rawData,
		KeyValues: keyValues,
	}, nil
}

func flattenJSON(prefix string, data interface{}, result *[]KeyValue) {
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			switch child := val.(type) {
			case map[string]interface{}:
				flattenJSON(key, child, result)
			case []interface{}:
				for i, item := range child {
					arrayKey := fmt.Sprintf("%s[%d]", key, i)
					flattenJSON(arrayKey, item, result)
				}
			default:
				*result = append(*result, KeyValue{
					Key:   key,
					Value: val,
				})
			}
		}
	default:
		if prefix != "" {
			*result = append(*result, KeyValue{
				Key:   prefix,
				Value: v,
			})
		}
	}
}
