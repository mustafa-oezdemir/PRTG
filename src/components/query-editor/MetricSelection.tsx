import React from 'react';
import { SelectableValue } from '@grafana/data';
import { AsyncMultiSelect, Combobox, InlineField, Stack, type ComboboxOption } from '@grafana/ui';
import { MyQuery, queryTypeOptions } from '../../types';

interface MetricSelectionProps {
  query: MyQuery;
  isLoading: boolean;
  deviceLoading: boolean;
  sensorLoading: boolean;
  groupOptions: Array<ComboboxOption<string>>;
  deviceOptions: Array<ComboboxOption<string>>;
  sensorOptions: Array<ComboboxOption<string>>;
  selectedGroup: ComboboxOption<string> | null;
  selectedDevice: ComboboxOption<string> | null;
  selectedSensor: ComboboxOption<string> | null;
  sensorId: string;
  channelQuery: string[];
  loadChannelOptions: () => Promise<Array<SelectableValue<string>>>;
  onQueryTypeChange: (option: ComboboxOption<string> | null) => void;
  onGroupChange: (option: ComboboxOption<string> | null) => void;
  onDeviceChange: (option: ComboboxOption<string> | null) => void;
  onSensorChange: (option: ComboboxOption<string> | null) => void;
  onChannelChange: (values: Array<SelectableValue<string>>) => void;
}

export function MetricSelection({
  query,
  isLoading,
  deviceLoading,
  sensorLoading,
  groupOptions,
  deviceOptions,
  sensorOptions,
  selectedGroup,
  selectedDevice,
  selectedSensor,
  sensorId,
  channelQuery,
  loadChannelOptions,
  onQueryTypeChange,
  onGroupChange,
  onDeviceChange,
  onSensorChange,
  onChannelChange,
}: MetricSelectionProps) {
  return (
    <Stack direction="row" gap={2}>
      <Stack direction="column" gap={1}>
        <InlineField label="Query Type" labelWidth={20} grow>
          <Combobox
            id="query-editor-queryType"
            options={queryTypeOptions}
            value={query.queryType}
            onChange={onQueryTypeChange}
            width={47}
          />
        </InlineField>

        <InlineField label="Group" labelWidth={20} grow>
          <Combobox
            id="query-editor-group"
            loading={isLoading}
            options={groupOptions}
            value={selectedGroup}
            onChange={onGroupChange}
            width={47}
            createCustomValue={true}
            isClearable={true}
            invalid={!query.queryType}
            placeholder="Select Group or type '*'"
          />
        </InlineField>

        <InlineField label="Device" labelWidth={20} grow>
          <Combobox
            id="query-editor-device"
            loading={deviceLoading}
            options={deviceOptions}
            value={selectedDevice}
            onChange={onDeviceChange}
            width={47}
            createCustomValue={true}
            isClearable={true}
            invalid={!query.group}
            placeholder="Select Device or type '*'"
          />
        </InlineField>
      </Stack>

      <Stack direction="column" gap={1}>
        <InlineField label="Sensor" labelWidth={20} grow>
          <Combobox
            id="query-editor-sensor"
            loading={sensorLoading}
            options={sensorOptions}
            value={selectedSensor}
            onChange={onSensorChange}
            width={47}
            createCustomValue={true}
            isClearable={true}
            invalid={!query.device}
            placeholder="Select Sensor or type '*'"
          />
        </InlineField>

        <InlineField label="Channel" labelWidth={20} grow>
          <AsyncMultiSelect
            id="query-editor-channel"
            key={sensorId}
            loadOptions={loadChannelOptions}
            defaultOptions={true}
            value={channelQuery.map((channel) => ({
              label: channel,
              value: channel,
            }))}
            onChange={onChannelChange}
            width={47}
            placeholder={sensorId ? 'Select Channels (multiple allowed)' : 'First select a sensor'}
            isClearable
            isDisabled={!sensorId}
            noOptionsMessage="No channels available"
          />
        </InlineField>
      </Stack>
    </Stack>
  );
}
