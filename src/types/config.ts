import { DataSourceJsonData } from '@grafana/data';

export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
  cacheTime?: number;
  timeZone?: string;
}

export interface MySecureJsonData {
  apiKey?: string;
}
