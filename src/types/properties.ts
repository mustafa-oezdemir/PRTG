export const filterPropertyList = [
  { name: 'active', visible_name: 'Active' },
  { name: 'message_raw', visible_name: 'Message' },
  { name: 'priority', visible_name: 'Priority' },
  { name: 'status', visible_name: 'Status' },
  { name: 'tags', visible_name: 'Tags' },
] as const;

export type FilterPropertyItem = typeof filterPropertyList[number];

export interface FilterPropertyOption {
  label: string;
  value: FilterPropertyItem['name'];
}

export const propertyList = [
  { name: 'group', visible_name: 'Group' },
  { name: 'device', visible_name: 'Device' },
  { name: 'sensor', visible_name: 'Sensor' },
] as const;

export type PropertyItem = typeof propertyList[number];

export interface PropertyOption {
  label: string;
  value: PropertyItem['name'];
}
