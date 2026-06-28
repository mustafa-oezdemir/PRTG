import { DataQuery } from '@grafana/schema';
import { StreamingConfig } from './streaming';

export enum QueryType {
  Metrics = 'metrics',
  Raw = 'raw',
  Text = 'text',
  Manual = 'manual',
}

export interface MyQuery extends DataQuery {
  queryType: QueryType;
  group: string;
  groupId: string;
  device: string;
  deviceId: string;
  sensor: string;
  sensorId: string;
  channel: string;
  channelArray: string[];
  manualMethod?: string;
  manualObjectId?: string;
  property?: string;
  filterProperty?: string;
  includeGroupName?: boolean;
  includeDeviceName?: boolean;
  includeSensorName?: boolean;
  refId: string;
  streaming?: StreamingConfig;
  isStreaming?: boolean;
  streamInterval?: number;
  streamId?: string;
  panelId?: string | number;
  queryId?: string;
  cacheTime?: number;
  bufferSize?: number;
  updateMode?: 'full' | 'append';
}

export interface QueryTypeOptions {
  label: string;
  value: QueryType;
}

export const queryTypeOptions = Object.keys(QueryType).map((key) => ({
  label: key,
  value: QueryType[key as keyof typeof QueryType],
}));
