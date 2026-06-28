import { MetricFindValue } from '@grafana/data';
import { getTemplateSrv } from '@grafana/runtime';
import {
  PRTGChannelListResponse,
  PRTGDeviceListResponse,
  PRTGGroupListResponse,
  PRTGSensorListResponse,
} from '../types';

export interface PrtgResourceClient {
  getGroups(): Promise<PRTGGroupListResponse>;
  getDevices(group: string): Promise<PRTGDeviceListResponse>;
  getSensors(device: string): Promise<PRTGSensorListResponse>;
  getChannels(sensorId: string): Promise<PRTGChannelListResponse>;
}

export async function runMetricFindQuery(
  client: PrtgResourceClient,
  query: string | { query?: string }
): Promise<MetricFindValue[]> {
  const queryText = typeof query === 'string' ? query : query.query || '';
  const [kind, ...argParts] = queryText.split(':');
  const rawArg = argParts.join(':');
  const arg = getTemplateSrv().replace(rawArg).trim();

  switch (kind.trim()) {
    case 'groups': {
      const response = await client.getGroups();
      return response.groups
        .filter((item) => item.group)
        .map((item) => ({ text: item.group, value: item.group }));
    }
    case 'devices': {
      if (!arg) {
        return [];
      }

      const response = await client.getDevices(arg);
      return response.devices
        .filter((item) => item.device)
        .map((item) => ({ text: item.device, value: item.device }));
    }
    case 'sensors': {
      if (!arg) {
        return [];
      }

      const response = await client.getSensors(arg);
      return response.sensors
        .filter((item) => item.sensor && item.objid)
        .map((item) => ({ text: item.sensor, value: String(item.objid) }));
    }
    case 'channels': {
      if (!arg) {
        return [];
      }

      const response = await client.getChannels(arg);
      const channelData = response.values?.[0];
      if (!channelData) {
        return [];
      }

      return Object.keys(channelData)
        .filter((key) => key !== 'datetime')
        .map((key) => ({ text: key, value: key }));
    }
    default:
      return [];
  }
}
