export interface ListItem {
  name: string;
  visible_name: string;
}

export interface PRTGItem {
  active: boolean;
  active_raw: number;
  channel: string;
  channel_raw: string;
  datetime: string;
  datetime_raw: number;
  device: string;
  device_raw: string;
  group: string;
  group_raw: string;
  message: string;
  message_raw: string;
  objid: number;
  objid_raw: number;
  priority: string;
  priority_raw: number;
  sensor: string;
  sensor_raw: string;
  status: string;
  status_raw: number;
  tags: string;
  tags_raw: string;
}

export interface PRTGGroupListResponse {
  prtgversion: string;
  treesize: number;
  groups: PRTGItem[];
}

export interface PRTGGroupResponse {
  groups: PRTGItem[];
}

export interface PRTGDeviceListResponse {
  prtgversion: string;
  treesize: number;
  devices: PRTGItem[];
}

export interface PRTGDeviceResponse {
  devices: PRTGItem[];
}

export interface PRTGSensorListResponse {
  prtgversion: string;
  treesize: number;
  sensors: PRTGItem[];
}

export interface PRTGSensorResponse {
  sensors: PRTGItem[];
}

export interface PRTGChannelListResponse {
  prtgversion: string;
  treesize: number;
  values: PRTGItemChannel[];
}

export interface PRTGItemChannel {
  [key: string]: number | string;
  datetime: string;
}
