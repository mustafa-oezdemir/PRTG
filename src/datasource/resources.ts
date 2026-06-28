import {
  PRTGChannelListResponse,
  PRTGDeviceListResponse,
  PRTGGroupListResponse,
  PRTGSensorListResponse,
} from '../types';

export type ResourceGetter = <T = unknown>(path: string) => Promise<T>;

export function getGroupsResource(getResource: ResourceGetter): Promise<PRTGGroupListResponse> {
  return getResource<PRTGGroupListResponse>('groups');
}

export function getDevicesResource(getResource: ResourceGetter, group: string): Promise<PRTGDeviceListResponse> {
  if (!group) {
    throw new Error('group is required');
  }

  return getResource<PRTGDeviceListResponse>(`devices/${encodeURIComponent(group)}`);
}

export function getSensorsResource(getResource: ResourceGetter, device: string): Promise<PRTGSensorListResponse> {
  if (!device) {
    throw new Error('device is required');
  }

  return getResource<PRTGSensorListResponse>(`sensors/${encodeURIComponent(device)}`);
}

export function getChannelsResource(getResource: ResourceGetter, sensorId: string): Promise<PRTGChannelListResponse> {
  if (!sensorId) {
    throw new Error('sensorId is required');
  }

  return getResource<PRTGChannelListResponse>(`channels/${encodeURIComponent(sensorId)}`);
}
