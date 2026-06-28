package prtg

type AnnotationQuery struct {
	From         int64    `json:"from,omitempty"`
	To           int64    `json:"to,omitempty"`
	Limit        int64    `json:"limit,omitempty"`
	AlertID      int64    `json:"alertId,omitempty"`
	DashboardID  int64    `json:"dashboardId,omitempty"`
	DashboardUID string   `json:"dashboardUID,omitempty"`
	PanelID      int64    `json:"panelId,omitempty"`
	UserID       int64    `json:"userId,omitempty"`
	Type         string   `json:"type,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	SensorID     string   `json:"sensorId,omitempty"`
}

type Annotation struct {
	ID      string                 `json:"id"`
	Time    int64                  `json:"time"`
	TimeEnd int64                  `json:"timeEnd"`
	Title   string                 `json:"title"`
	Text    string                 `json:"text"`
	Tags    []string               `json:"tags"`
	Type    string                 `json:"type,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type AnnotationResponse struct {
	Annotations []Annotation `json:"annotations"`
	Total       int          `json:"total"`
}

type PrtgAnnotationResponse struct {
	Annotations []PrtgAnnotation `json:"annotations"`
}

type PrtgAnnotation struct {
	ID      int64    `json:"id"`
	Time    int64    `json:"time"`
	TimeEnd int64    `json:"timeEnd"`
	Text    string   `json:"text"`
	Tags    []string `json:"tags"`
}
