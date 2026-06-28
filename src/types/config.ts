import { DataSourceJsonData } from '@grafana/data';

export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
  cacheTime?: number;
  timeZone?: string;
  timezone?: string;
  selectedTimezone?: string;
}

export interface MySecureJsonData {
  apiKey?: string;
}
