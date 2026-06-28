package prtg

import "time"

type PRTGAPI interface {
	GetGroups() (*PrtgGroupListResponse, error)
	GetStatusList() (*PrtgStatusListResponse, error)
	GetDevices(groupId string) (*PrtgDevicesListResponse, error)
	GetSensors(deviceId string) (*PrtgSensorsListResponse, error)
	GetChannels(sensorId string) (*PrtgChannelValueStruct, error)
	GetHistoricalData(sensorId string, from time.Time, to time.Time) (*PrtgHistoricalDataResponse, error)
	ExecuteManualMethod(method string, objectId string) (*PrtgManualMethodResponse, error)
	GetAnnotationData(query *AnnotationQuery) (*AnnotationResponse, error)
	GetCacheTime() time.Duration
}

type ApiInterface interface {
	GetCacheTime() time.Duration
	SetTimeout(timeout time.Duration)
	GetStatusList() (*PrtgStatusListResponse, error)
	GetGroups() (*PrtgGroupListResponse, error)
	GetDevices(group string) (*PrtgDevicesListResponse, error)
	GetSensors(device string) (*PrtgSensorsListResponse, error)
	GetChannels(sensorId string) (*PrtgChannelValueStruct, error)
	GetHistoricalData(sensorID string, startDate, endDate time.Time) (*PrtgHistoricalDataResponse, error)
	ExecuteManualMethod(method string, objectId string) (*PrtgManualMethodResponse, error)
	GetAnnotationData(query *AnnotationQuery) (*AnnotationResponse, error)
}
