import { Dispatch, SetStateAction, useEffect, useMemo, useState } from 'react';
import { SelectableValue } from '@grafana/data';
import type { ComboboxOption } from '@grafana/ui';
import { DataSource } from '../../../datasource';
import { MyQuery, filterPropertyList, propertyList } from '../../../types';

export interface QueryEditorLists {
  groups: Array<ComboboxOption<string>>;
  devices: Array<ComboboxOption<string>>;
  sensors: Array<ComboboxOption<string>>;
  channels: Array<ComboboxOption<string>>;
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
  return [...items].sort((a, b) =>
    (a.label ?? '').localeCompare(b.label ?? '')
  );
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

  /**
   * Load PRTG groups.
   */
  useEffect(() => {
    let cancelled = false;

    async function fetchGroups() {
      setIsLoading(true);

      try {
        const response = await datasource.getGroups();

        if (cancelled) {
          return;
        }

        if (response && Array.isArray(response.groups)) {
          const groupOptions: Array<ComboboxOption<string>> =
            response.groups.map((groupItem) => ({
              label: groupItem.group,
              value: groupItem.group.toString(),
            }));

          setLists((prev) => ({
            ...prev,
            groups: groupOptions,
          }));
        } else {
          console.error('Invalid groups response format:', response);

          setLists((prev) => ({
            ...prev,
            groups: [],
          }));
        }
      } catch (error) {
        if (cancelled) {
          return;
        }

        console.error('Error fetching groups:', error);

        setLists((prev) => ({
          ...prev,
          groups: [],
        }));
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    }

    fetchGroups();

    return () => {
      cancelled = true;
    };
  }, [datasource]);

  /**
   * Load devices after a group has been selected.
   */
  useEffect(() => {
    let cancelled = false;

    async function fetchDevices() {
      if (!group) {
        setLists((prev) => ({
          ...prev,
          devices: [],
          sensors: [],
          channels: [],
        }));

        return;
      }

      setIsLoading(true);

      try {
        const response = await datasource.getDevices(group);

        if (cancelled) {
          return;
        }

        if (response && Array.isArray(response.devices)) {
          const filteredDevices = response.devices.filter(
            (deviceItem) => deviceItem.group === group
          );

          const deviceOptions: Array<ComboboxOption<string>> =
            filteredDevices.map((deviceItem) => ({
              label: deviceItem.device,
              value: deviceItem.device.toString(),
            }));

          setLists((prev) => ({
            ...prev,
            devices: deviceOptions,
            sensors: [],
            channels: [],
          }));
        } else {
          console.error('Invalid devices response format:', response);

          setLists((prev) => ({
            ...prev,
            devices: [],
            sensors: [],
            channels: [],
          }));
        }
      } catch (error) {
        if (cancelled) {
          return;
        }

        console.error('Error fetching devices:', error);

        setLists((prev) => ({
          ...prev,
          devices: [],
          sensors: [],
          channels: [],
        }));
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    }

    fetchDevices();

    return () => {
      cancelled = true;
    };
  }, [datasource, group]);

  /**
   * Load sensors after a device has been selected.
   */
  useEffect(() => {
    let cancelled = false;

    async function fetchSensors() {
      if (!device) {
        setLists((prev) => ({
          ...prev,
          sensors: [],
          channels: [],
        }));

        return;
      }

      setIsLoading(true);

      try {
        const response = await datasource.getSensors(device);

        if (cancelled) {
          return;
        }

        if (response && Array.isArray(response.sensors)) {
          const filteredSensors = response.sensors.filter(
            (sensorItem) => sensorItem.device === device
          );

          const sensorOptions: Array<ComboboxOption<string>> =
            filteredSensors.map((sensorItem) => ({
              label: sensorItem.sensor,
              value: sensorItem.sensor.toString(),
            }));

          setLists((prev) => ({
            ...prev,
            sensors: sensorOptions,
            channels: [],
          }));
        } else {
          console.error('Invalid sensors response format:', response);

          setLists((prev) => ({
            ...prev,
            sensors: [],
            channels: [],
          }));
        }
      } catch (error) {
        if (cancelled) {
          return;
        }

        console.error('Error fetching sensors:', error);

        setLists((prev) => ({
          ...prev,
          sensors: [],
          channels: [],
        }));
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    }

    fetchSensors();

    return () => {
      cancelled = true;
    };
  }, [datasource, device]);

  /**
   * Load channels after a sensor has been selected.
   */
  useEffect(() => {
    let cancelled = false;

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

        if (cancelled) {
          return;
        }

        if (!response) {
          console.error('Empty channels response received');

          setLists((prev) => ({
            ...prev,
            channels: [],
          }));

          return;
        }

        if (
          response.values &&
          Array.isArray(response.values) &&
          response.values.length > 0
        ) {
          const channelData = (response.values[0] || {}) as Record<
            string,
            unknown
          >;

          const channelOptions: Array<ComboboxOption<string>> =
            Object.entries(channelData)
              .filter(([key]) => key !== 'datetime')
              .map(([key]) => ({
                label: key,
                value: key,
              }));

          setLists((prev) => ({
            ...prev,
            channels: channelOptions,
          }));

          if (
            query.channel &&
            channelOptions.some(
              (option) => option.value === query.channel
            )
          ) {
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
        if (cancelled) {
          return;
        }

        console.error('Error fetching channels:', error);

        setLists((prev) => ({
          ...prev,
          channels: [],
        }));
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    }

    fetchChannels();

    return () => {
      cancelled = true;
    };
  }, [datasource, sensorId, query.channel, setChannel]);

  /**
   * Load property lists for Text and Raw modes.
   */
  useEffect(() => {
    if (isTextMode || isRawMode) {
      const propertyOptions: Array<SelectableValue<string>> =
        propertyList.map((item) => ({
          label: item.visible_name,
          value: item.name,
        }));

      const filterPropertyOptions: Array<SelectableValue<string>> =
        filterPropertyList.map((item) => ({
          label: item.visible_name,
          value: item.name,
        }));

      setLists((prev) => ({
        ...prev,
        properties: propertyOptions,
        filterProperties: filterPropertyOptions,
      }));
    } else {
      setLists((prev) => ({
        ...prev,
        properties: [],
        filterProperties: [],
      }));
    }
  }, [isTextMode, isRawMode]);

  /**
   * Sort selectable options alphabetically.
   */
  const groupOptions = useMemo(
    () => sortByLabel(lists.groups),
    [lists.groups]
  );

  const deviceOptions = useMemo(
    () => sortByLabel(lists.devices),
    [lists.devices]
  );

  const sensorOptions = useMemo(
    () => sortByLabel(lists.sensors),
    [lists.sensors]
  );

  /**
   * Resolve selected values.
   *
   * The fallback objects are important when a saved query contains
   * a value before the corresponding async options have been loaded.
   */
  const selectedGroup = useMemo(() => {
    return (
      groupOptions.find((option) => option.value === group) ||
      (group
        ? {
            label: group,
            value: group,
          }
        : null)
    );
  }, [groupOptions, group]);

  const selectedDevice = useMemo(() => {
    return (
      deviceOptions.find((option) => option.value === device) ||
      (device
        ? {
            label: device,
            value: device,
          }
        : null)
    );
  }, [deviceOptions, device]);

  const selectedSensor = useMemo(() => {
    return (
      sensorOptions.find((option) => option.value === sensor) ||
      (sensor
        ? {
            label: sensor,
            value: sensor,
          }
        : null)
    );
  }, [sensorOptions, sensor]);

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
  };
}
