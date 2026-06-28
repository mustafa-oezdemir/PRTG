import { 
  DataSourceInstanceSettings, 
  ScopedVars, 
  DataQueryRequest,
  DataQueryResponse,
  MetricFindValue,
} from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { Observable } from 'rxjs';
import { processAnnotationEvents } from './datasource/annotations';
import { runMetricFindQuery } from './datasource/metricFind';
import { filterPrtgQuery, runPrtgDataQuery } from './datasource/querying';
import {
  getChannelsResource,
  getDevicesResource,
  getGroupsResource,
  getSensorsResource,
} from './datasource/resources';
import { applyQueryTemplateVariables } from './datasource/templating';
import {
  MyQuery,
  MyDataSourceOptions,
  PRTGGroupListResponse,
  PRTGDeviceListResponse,
  PRTGSensorListResponse,
  PRTGChannelListResponse,
} from './types'

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {

  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: MyQuery, scopedVars: ScopedVars) {
    return applyQueryTemplateVariables(query, scopedVars);
  }

  filterQuery(query: MyQuery): boolean {
    return filterPrtgQuery(query);
  }

  async metricFindQuery(query: string | { query?: string }): Promise<MetricFindValue[]> {
    return runMetricFindQuery(this, query);
  }

  async getGroups(): Promise<PRTGGroupListResponse> {
    return getGroupsResource((path) => this.getResource(path))
  }

  async getDevices(group: string): Promise<PRTGDeviceListResponse> {
    return getDevicesResource((path) => this.getResource(path), group)
  }

  async getSensors(device: string): Promise<PRTGSensorListResponse> {
    return getSensorsResource((path) => this.getResource(path), device);
  }

  async getChannels(sensorId: string): Promise<PRTGChannelListResponse> {
    return getChannelsResource((path) => this.getResource(path), sensorId);
  }

  annotations = {
    QueryEditor: undefined,
    processEvents: processAnnotationEvents,
  };

  query(request: DataQueryRequest<MyQuery>): Observable<DataQueryResponse> {
    return runPrtgDataQuery(this.uid, request, (backendRequest) => super.query(backendRequest));
  }

  // Stream control methods
  async getStreamStatus(streamId: string): Promise<any> {
    return this.getResource(`stream-status/${streamId}`);
  }

  async stopStream(streamId: string): Promise<void> {
    return this.getResource(`stop-stream/${streamId}`);
  }
}
