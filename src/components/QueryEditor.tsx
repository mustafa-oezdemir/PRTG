import React, { useEffect, useState, useCallback, ChangeEvent, useRef } from 'react';
import { Stack, type ComboboxOption } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data'
import { DataSource } from '../datasource'
import {
  DEFAULT_STREAMING_OPTIONS, MyDataSourceOptions, MyQuery, QueryType, manualApiMethods
} from '../types'
import { DisplayOptions } from './query-editor/DisplayOptions';
import { ManualQueryOptions } from './query-editor/ManualQueryOptions';
import { MetricSelection } from './query-editor/MetricSelection';
import { PropertyOptions } from './query-editor/PropertyOptions';
import { StreamingOptions } from './query-editor/StreamingOptions';
import { usePrtgSelectionLists } from './query-editor/hooks/usePrtgSelectionLists';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>

export function QueryEditor({ query, onChange, onRunQuery, datasource }: Props) {
  const prevQueryRef = useRef<MyQuery | null>(null);
  const runQueryIfChanged = useCallback(() => {
    const currentQuery = JSON.stringify({ ...query, refId: query.refId }); // Include refId in comparison
    const prevQuery = JSON.stringify(prevQueryRef.current);

    if (currentQuery !== prevQuery) {
      prevQueryRef.current = query;
      onRunQuery();
    }
  }, [query, onRunQuery]);

  const isMetricsMode = query.queryType === QueryType.Metrics
  const isRawMode = query.queryType === QueryType.Raw
  const isTextMode = query.queryType === QueryType.Text
  const isManualMode = query.queryType === QueryType.Manual

  /* ===================================================== HOOKS ============================================================*/
  const [group, setGroup] = useState<string>(query.group || '')
  const [device, setDevice] = useState<string>(query.device || '')
  const [sensor, setSensor] = useState<string>(query.sensor || '')
  const [, setChannel] = useState<string>(query.channel || '')
  const [channelQuery, setChannelQuery] = useState<string[]>(query.channelArray || [])
  const [sensorId, setSensorId] = useState<string>(query.sensorId || '')
  const [manualMethod, setManualMethod] = useState<string>(query.manualMethod || '');
  const [manualObjectId, setManualObjectId] = useState<string>(query.manualObjectId || '');
  const [streamIntervalValue, setStreamIntervalValue] = useState<string>(String(query.streamInterval || 2500));

  const {
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
  } = usePrtgSelectionLists({
    datasource,
    query,
    isTextMode,
    isRawMode,
    group,
    device,
    sensor,
    sensorId,
    setChannel,
  });
  /* ==================================================  INITIAL VALUES  ================================================== */
  useEffect(() => {
    setGroup((prev) => query.group ?? prev);
    setDevice((prev) => query.device ?? prev);
    setSensor((prev) => query.sensor ?? prev);
    setChannel((prev) => query.channel ?? prev);
    setSensorId((prev) => query.sensorId ?? prev);
    setManualMethod((prev) => query.manualMethod ?? prev);
    setManualObjectId((prev) => query.manualObjectId ?? prev);
    setStreamIntervalValue(String(query.streamInterval || 2500));
    // Add this line to restore channel selections
    setChannelQuery((prev) => query.channelArray || prev || []);
  }, [query]);


  /* ==================================================  FIND IDs ================================================= */
  const findGroupId = useCallback(async (groupName: string) => {
    try {
      const response = await datasource.getGroups()
      if (response && Array.isArray(response.groups)) {
        const group = response.groups.find((g) => g.group === groupName)
        if (group) {
          return group.objid.toString()
        }
      }
    } catch (error) {
      console.error('Error finding group ID:', error)
    }
    return ''
  }, [datasource])

  const findDeviceId = useCallback(async (deviceName: string) => {
    try {
      const response = await datasource.getDevices(group)
      if (response && Array.isArray(response.devices)) {
        const device = response.devices.find((d) => d.device === deviceName)
        if (device) {
          return device.objid.toString()
        }
      }
    } catch (error) {
      console.error('Error finding device ID:', error)
    }
    return ''
  }, [datasource, group])

  const findSensorObjid = useCallback(async (sensorName: string) => {
    try {
      const response = await datasource.getSensors(device)
      if (response && Array.isArray(response.sensors)) {
        const sensor = response.sensors.find((s) => s.sensor === sensorName)
        if (sensor) {
          setSensorId(sensor.objid.toString())
          return sensor.objid.toString()
        } else {
          console.error('Sensor not found:', sensorName)
        }
      } else {
        console.error('Invalid response format:', response)
      }
    } catch (error) {
      console.error('Error fetching sensors:', error)
    }
    return ''
  }, [datasource, device, setSensorId])
  /* ==================================================  EVENT HANDLERS ==================================================  */

  /* ==================================================  QUERY  ==================================================  */
  /* ==================================================  ONQUERYTYPESCHANGE ==================================================  */
  const onQueryTypeChange = useCallback((option: ComboboxOption<string> | null) => {
    if (option?.value) {
      onChange({ ...query, queryType: option.value as QueryType });
      runQueryIfChanged();
    }
  }, [query, onChange, runQueryIfChanged]);

  /* ==================================================  ONGROUPCHANGE ==================================================  */
  const onGroupChange = useCallback(async (option: ComboboxOption<string> | null) => {
    if (!option?.value) {
      return;
    }

    const groupObjId = await findGroupId(option.value);
    setGroup(option.value);

    const updatedQuery = {
      ...query,
      group: option.value,
      groupId: groupObjId,
    };
    onChange(updatedQuery);
    setLists(prev => ({ ...prev, devices: [], sensors: [], channels: [] }));
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged, findGroupId, setLists]);

  /* ==================================================  ONDEVICECHANGE ================================================= */
  const onDeviceChange = useCallback(async (option: ComboboxOption<string> | null) => {
    if (!option?.value) {
      return;
    }

    const deviceObjId = await findDeviceId(option.value);

    setDevice(option.value);
    const updatedQuery = {
      ...query,
      device: option.value,
      deviceId: deviceObjId,
    };
    onChange(updatedQuery);
    setLists(prev => ({ ...prev, sensors: [], channels: [] }));
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged, findDeviceId, setLists]);
  /* ==================================================  ONSENSORCHANGE ==================================================  */
  const onSensorChange = useCallback(async (option: ComboboxOption<string> | null) => {
    if (!option?.value) {
      return;
    }

    const sensorObjId = await findSensorObjid(option.value);

    setSensor(option.value);
    setSensorId(sensorObjId);
    setLists(prev => ({ ...prev, channels: [] }));

    const updatedQuery = {
      ...query,
      sensor: option.value,
      sensorId: sensorObjId,
    };
    onChange(updatedQuery);

    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged, findSensorObjid, setLists]);  /* ==================================================  ONCHANNELCHANGE ==================================================  */
  const onChannelChange = useCallback((values: Array<SelectableValue<string>>) => {
    const selectedChannels = values.map(v => v.value!);

    // Update local state
    setChannelQuery(selectedChannels);

    // CRITICAL: Update query to include ALL selected channels in a SINGLE query
    // This prevents Grafana from creating multiple queries (refId A, B, C...)
    const updatedQuery = {
      ...query,
      channel: selectedChannels[0] || '', // First channel for backward compatibility
      channelArray: selectedChannels, // ALL selected channels in one array
      // Generate series names for each channel
      seriesNames: selectedChannels.map(channel =>
        `${query.sensor || 'Sensor'} - ${channel}`
      ),
    };

    onChange(updatedQuery);    // Only trigger query execution if we have channels selected
    // Don't use runQueryIfChanged() as it might create duplicate queries
    if (selectedChannels.length > 0) {
      // Use timeout to ensure state is updated before running query
      setTimeout(() => {
        onRunQuery();
      }, 0);
    }
  }, [query, onChange, onRunQuery]);
  /* ==================================================  ON INCLUDE GROUP NAME ==================================================  */
  const onIncludeGroupName = (event: ChangeEvent<HTMLInputElement>) => {
    const updatedQuery = { ...query, includeGroupName: event.currentTarget.checked };
    onChange(updatedQuery);
    runQueryIfChanged();
  }

  /* ==================================================  ON INCLUDE DEVICE NAME ==================================================  */
  const onIncludeDeviceName = (event: React.ChangeEvent<HTMLInputElement>) => {
    const updatedQuery = { ...query, includeDeviceName: event.currentTarget.checked };
    onChange(updatedQuery);
    runQueryIfChanged();
  }

  /* ==================================================  ON INCLUDE SENSOR NAME ==================================================  */
  const onIncludeSensorName = (event: ChangeEvent<HTMLInputElement>) => {
    const updatedQuery = { ...query, includeSensorName: event.currentTarget.checked };
    onChange(updatedQuery);
    runQueryIfChanged();
  }  /* ==================================================  ON MANUAL OBJECT ID CHANGE ==================================================  */
  const onManualObjectIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const value = event.currentTarget.value;
    setManualObjectId(value);
    const updatedQuery = {
      ...query,
      manualObjectId: value,
    };
    onChange(updatedQuery);
    runQueryIfChanged();
  };

  const onManualMethodChange = useCallback((value: string) => {
    setManualMethod(value);
    const updatedQuery = {
      ...query,
      manualMethod: value,
    };
    onChange(updatedQuery);
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged]);

  const onPropertyChange = useCallback((value: string) => {
    const updatedQuery = { ...query, property: value };
    onChange(updatedQuery);
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged]);

  const onFilterPropertyChange = useCallback((value: string) => {
    const updatedQuery = { ...query, filterProperty: value };
    onChange(updatedQuery);
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged]);

  /* ==================================================  STREAM INTERVAL HANDLERS ==================================================  */
  const handleStreamIntervalChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    const value = e.currentTarget.value;
    setStreamIntervalValue(value);
  }, []);

  const handleStreamIntervalBlur = useCallback(() => {
    const interval = Math.max(0, Math.min(60000, parseInt(streamIntervalValue, 10) || 2500));
    const updatedQuery = {
      ...query,
      streamInterval: interval,
    };
    onChange(updatedQuery);
    runQueryIfChanged();
  }, [streamIntervalValue, query, onChange, runQueryIfChanged]);

  const onStreamingToggle = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    const isStreaming = event.currentTarget.checked;
    const streamInterval = isStreaming ? (query.streamInterval || DEFAULT_STREAMING_OPTIONS.defaultInterval) : undefined;
    const updatedQuery = {
      ...query,
      isStreaming,
      streamInterval,
      updateMode: isStreaming ? DEFAULT_STREAMING_OPTIONS.updateMode : query.updateMode,
      bufferSize: isStreaming ? (query.bufferSize || DEFAULT_STREAMING_OPTIONS.bufferSize) : query.bufferSize,
      cacheTime: isStreaming ? (query.cacheTime || DEFAULT_STREAMING_OPTIONS.cacheTime) : query.cacheTime,
    };
    onChange(updatedQuery);
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged]);



  /* ================================================== DESTRUCTURING ================================================== */
  // Set default streaming values
  useEffect(() => {
    if (query.isStreaming === undefined) {
      const updatedQuery = {
        ...query,
        isStreaming: false,
        streamInterval: 2500, // Default interval 5ms (2,5 seconds)
      };
      onChange(updatedQuery);
    }
  }, [query, onChange]);

  /* ================================================== RENDER ================================================== */
  return (
    <Stack direction="column" gap={2}>
      <MetricSelection
        query={query}
        isLoading={isLoading}
        deviceLoading={!lists.devices.length && !!query.group}
        sensorLoading={!lists.sensors.length && !!query.device}
        groupOptions={groupOptions}
        deviceOptions={deviceOptions}
        sensorOptions={sensorOptions}
        selectedGroup={selectedGroup}
        selectedDevice={selectedDevice}
        selectedSensor={selectedSensor}
        sensorId={sensorId}
        channelQuery={channelQuery}
        loadChannelOptions={loadChannelOptions}
        onQueryTypeChange={onQueryTypeChange}
        onGroupChange={onGroupChange}
        onDeviceChange={onDeviceChange}
        onSensorChange={onSensorChange}
        onChannelChange={onChannelChange}
      />


      {/* Show display name options for both Metrics and Streaming */}
      {(isMetricsMode || query.isStreaming || isRawMode || isTextMode) && (
        <DisplayOptions
          query={query}
          onIncludeGroupName={onIncludeGroupName}
          onIncludeDeviceName={onIncludeDeviceName}
          onIncludeSensorName={onIncludeSensorName}
        />
      )}
      
      {/* Options for Text and Raw modes */}
      {(isTextMode || isRawMode) && (
        <PropertyOptions
          property={query.property}
          filterProperty={query.filterProperty}
          properties={lists.properties}
          filterProperties={lists.filterProperties}
          onPropertyChange={onPropertyChange}
          onFilterPropertyChange={onFilterPropertyChange}
        />
      )}      {/* Manual API Query Section */}
      {isManualMode && (
        <ManualQueryOptions
          manualMethods={manualApiMethods}
          manualMethod={manualMethod}
          manualObjectId={manualObjectId}
          sensorId={sensorId}
          onManualMethodChange={onManualMethodChange}
          onManualObjectIdChange={onManualObjectIdChange}
        />
      )}

      {/* Always show streaming options */}
      <StreamingOptions
        query={query}
        streamIntervalValue={streamIntervalValue}
        onStreamingToggle={onStreamingToggle}
        onStreamIntervalChange={handleStreamIntervalChange}
        onStreamIntervalBlur={handleStreamIntervalBlur}
      />

    </Stack>
  )
}
