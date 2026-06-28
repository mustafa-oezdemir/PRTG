import { ScopedVars } from '@grafana/data';
import { getTemplateSrv } from '@grafana/runtime';
import { MyQuery } from '../types';

export function applyQueryTemplateVariables(query: MyQuery, scopedVars: ScopedVars): MyQuery {
  const templateSrv = getTemplateSrv();
  const replace = (value?: string) => templateSrv.replace(value || '', scopedVars);
  const channel = replace(query.channel);
  const channelArray = query.channelArray?.map((item) => replace(item)).filter(Boolean) || [];

  return {
    ...query,
    group: replace(query.group),
    groupId: replace(query.groupId),
    device: replace(query.device),
    deviceId: replace(query.deviceId),
    sensor: replace(query.sensor),
    sensorId: replace(query.sensorId),
    channel,
    channelArray,
    manualObjectId: replace(query.manualObjectId),
  };
}
