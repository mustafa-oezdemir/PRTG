package prtg

import (
	"fmt"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtgtime"
)

func (a *Api) GetAnnotationData(query *AnnotationQuery) (*AnnotationResponse, error) {
	fromTime := time.Unix(0, query.From*int64(time.Millisecond))
	toTime := time.Unix(0, query.To*int64(time.Millisecond))

	histData, err := a.GetHistoricalData(query.SensorID, fromTime, toTime)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical data for annotations: %w", err)
	}

	annotations := make([]Annotation, 0)
	for i, data := range histData.HistData {
		t, _, err := prtgtime.ParseDateTime(data.Datetime)
		if err != nil {
			continue
		}

		uid := fmt.Sprintf("uid:%s_%d", query.SensorID, i)

		annotation := Annotation{
			ID:      uid,
			Time:    t.UnixMilli(),
			TimeEnd: t.UnixMilli(),
			Title:   fmt.Sprintf("Sensor: %s", query.SensorID),
			Tags:    []string{"prtg", fmt.Sprintf("sensor:%s", query.SensorID)},
			Type:    "annotation",
			Data:    data.Value,
		}

		annotations = append(annotations, annotation)
	}

	if query.Limit > 0 && int64(len(annotations)) > query.Limit {
		annotations = annotations[:query.Limit]
	}

	return &AnnotationResponse{
		Annotations: annotations,
		Total:       len(annotations),
	}, nil
}
