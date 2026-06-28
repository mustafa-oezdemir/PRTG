import { Dispatch, SetStateAction, useEffect, useMemo, useState } from 'react';
import { SelectableValue } from '@grafana/data';
import type { ComboboxOption } from '@grafana/ui';
import { DataSource } from '../../../datasource';
import { MyQuery, filterPropertyList, propertyList } from '../../../types';

export interface QueryEditorLists {
  groups: Array<ComboboxOption<string>>;
  devices: Array<ComboboxOption<string>>;
  sensors: Array<ComboboxOption<string>>;
  channels: Array<SelectableValue<string>>;
  values: Array<SelectableValue<string>>;
  properties: Array<SelectableValue<string>>;
  filterProperties: Array<SelectableValue<string>>;
}

interface UsePrtgSelectionListsArgs {
  datasource: DataSource;
  query: MyQuery;
  isTextMode: boolean;
  isRawMode: boolean;
  group: string;
  device: string;
  sensor: string;
  sensorId: string;
  setChannel: Dispatch<SetStateAction<string>>;
}

const emptyLists: QueryEditorLists = {
  groups: [],
  devices: [],
  sensors: [],
  channels: [],
  values: [],
  properties: [],
  filterProperties: [],
};

function sortByLabel<T extends { label?: string }>(items: T[]): T[] {
  return [...items].sort((a, b) => (a.label ?? '').localeCompare(b.label ?? ''));
}

export function usePrtgSelectionLists({
  datasource,
  query,
  isTextMode,
  isRawMode,
  group,
  device,
  sensor,
  sensorId,
  setChannel,
}: UsePrtgSelectionListsArgs) {
  const [lists, setLists] = useState<QueryEditorLists>(emptyLists);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    async function fetchGroups() {
      setIsLoading(true);
      try {
        const response = await datasource.getGroups();
        if (response && Array.isArray(response.groups)) {
          const groupOptions = response.groups.map((group) => ({
            label: group.group,
            value: group.group.toString(),
          }));

          setTimeout(() => {
            setLists((prev) => ({
              ...prev,
              groups: groupOptions,
            }));
            setIsLoading(false);
          }, 0);
        } else {
          console.error('Invalid response format:', response);
          setTimeout(() => {
            setLists((prev) => ({
              ...prev,
              groups: [],
            }));
            setIsLoading(false);
          }, 0);
        }
      } catch (error) {
        console.error('Error fetching groups:', error);
        setTimeout(() => {
          setLists((prev) => ({
            ...prev,
            groups: [],
          }));
          setIsLoading(false);
        }, 0);
      }
    }

    fetchGroups();
  }, [datasource]);

  useEffect(() => {
    async function fetchDevices() {
      if (!group) {
        return;
      }

      setIsLoading(true);
      try {
        const response = await datasource.getDevices(group);
        if (response && Array.isArray(response.devices)) {
          const filteredDevices = group ? response.devices.filter((device) => device.group === group) : response.devices;
          const deviceOptions = filteredDevices.map((device) => ({
            label: device.device,
            value: device.device.toString(),
          }));
          setLists((prev) => ({
            ...prev,
            devices: deviceOptions,
          }));
        } else {
          console.error('Invalid devices response format:', response);
          setLists((prev) => ({
            ...prev,
            devices: [],
          }));
        }
      } catch (error) {
        console.error('Error fetching devices:', error);
        setLists((prev) => ({
          ...prev,
          devices: [],
        }));
      } finally {
        setIsLoading(false);
      }
    }

    fetchDevices();
  }, [datasource, group]);

  useEffect(() => {
    async function fetchSensors() {
      if (!device) {
        return;
      }

      setIsLoading(true);
      try {
        const response = await datasource.getSensors(device);
        if (response && Array.isArray(response.sensors)) {
          const filteredSensors = device ? response.sensors.filter((sensor) => sensor.device === device) : response.sensors;
          const sensorOptions = filteredSensors.map((sensor) => ({
            label: sensor.sensor,
            value: sensor.sensor.toString(),
          }));
          setLists((prev) => ({
            ...prev,
            sensors: sensorOptions,
          }));
        } else {
          console.error('Invalid sensors response format:', response);
          setLists((prev) => ({
            ...prev,
            sensors: [],
          }));
        }
      } catch (error) {
        console.error('Error fetching sensors:', error);
        setLists((prev) => ({
          ...prev,
          sensors: [],
        }));
      } finally {
        setIsLoading(false);
      }
    }

    fetchSensors();
  }, [datasource, device]);

  useEffect(() => {
    async function fetchChannels() {
      if (!sensorId) {
        setLists((prev) => ({
          ...prev,
          channels: [],
        }));
        return;
      }

      setIsLoading(true);
      try {
        const response = await datasource.getChannels(sensorId);
        if (!response) {
          console.error('Empty response received');
          setLists((prev) => ({
            ...prev,
            channels: [],
          }));
          return;
        }

        if (response.values && Array.isArray(response.values) && response.values.length > 0) {
          const channelData = (response.values[0] || {}) as Record<string, unknown>;
          const channelOptions = Object.entries(channelData)
            .filter(([key]) => key !== 'datetime')
            .map(([key]) => ({
              label: key,
              value: key,
            }));

          setLists((prev) => ({
            ...prev,
            channels: channelOptions,
          }));

          if (query.channel && channelOptions.some((option) => option.value === query.channel)) {
            setChannel(query.channel);
          }
        } else {
          console.warn('No channel data found in response');
          setLists((prev) => ({
            ...prev,
            channels: [],
          }));
        }
      } catch (error) {
        console.error('Error fetching channels:', error);
        setLists((prev) => ({
          ...prev,
          channels: [],
        }));
      } finally {
        setIsLoading(false);
      }
    }

    fetchChannels();
  }, [datasource, sensorId, query.channel, setChannel]);

  useEffect(() => {
    if (isTextMode || isRawMode) {
      const propertyOptions = propertyList.map((item) => ({
        label: item.visible_name,
        value: item.name,
      }));
      const filterPropertyOptions = filterPropertyList.map((item) => ({
        label: item.visible_name,
        value: item.name,
      }));

      setLists((prev) => ({
        ...prev,
        properties: propertyOptions,
        filterProperties: filterPropertyOptions,
      }));
    }
  }, [isTextMode, isRawMode]);

  const groupOptions = useMemo(() => sortByLabel(lists.groups), [lists.groups]);
  const deviceOptions = useMemo(() => sortByLabel(lists.devices), [lists.devices]);
  const sensorOptions = useMemo(() => sortByLabel(lists.sensors), [lists.sensors]);

  const selectedGroup = useMemo(() => {
    return groupOptions.find((option) => option.value === group) || (group ? { label: group, value: group } : null);
  }, [groupOptions, group]);

  const selectedDevice = useMemo(() => {
    return deviceOptions.find((option) => option.value === device) || (device ? { label: device, value: device } : null);
  }, [deviceOptions, device]);

  const selectedSensor = useMemo(() => {
    return sensorOptions.find((option) => option.value === sensor) || (sensor ? { label: sensor, value: sensor } : null);
  }, [sensorOptions, sensor]);

  const loadChannelOptions = useMemo(() => async () => {
    if (!sensorId) {
      return [];
    }

    try {
      const response = await datasource.getChannels(sensorId);

      if (!response) {
        console.warn('No response received from getChannels');
        return [];
      }

      if (typeof response === 'object' && 'values' in response) {
        const values = response.values;
        if (!Array.isArray(values) || values.length === 0) {
          console.warn('No channel values found in response');
          return [];
        }

        const channelData = values[0];
        if (typeof channelData !== 'object') {
          console.warn('Invalid channel data format');
          return [];
        }

        return Object.keys(channelData)
          .filter((key) => key !== 'datetime')
          .map((key) => ({
            label: key,
            value: key,
          }));
      }

      console.warn('Unexpected response format:', response);
      return [];
    } catch (error: any) {
      console.error('Error loading channels:', error?.message || error);
      return [];
    }
  }, [sensorId, datasource]);

  return {
    lists,
    setLists,
    isLoading,
    groupOptions,
    deviceOptions,
    sensorOptions,
    selectedGroup,
    selectedDevice,
    selectedSensor,
    loadChannelOptions,
  };
}
