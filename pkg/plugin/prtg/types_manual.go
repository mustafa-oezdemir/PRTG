package prtg

type PrtgManualMethodResponse struct {
	Manuel    map[string]interface{} `json:"raw"`
	KeyValues []KeyValue             `json:"keyValues"`
}

type KeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
