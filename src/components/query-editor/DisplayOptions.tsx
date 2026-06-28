import React, { ChangeEvent } from 'react';
import { FieldSet, InlineField, InlineSwitch, Stack } from '@grafana/ui';
import { MyQuery } from '../../types';

interface DisplayOptionsProps {
  query: MyQuery;
  onIncludeGroupName: (event: ChangeEvent<HTMLInputElement>) => void;
  onIncludeDeviceName: (event: ChangeEvent<HTMLInputElement>) => void;
  onIncludeSensorName: (event: ChangeEvent<HTMLInputElement>) => void;
}

export function DisplayOptions({
  query,
  onIncludeGroupName,
  onIncludeDeviceName,
  onIncludeSensorName,
}: DisplayOptionsProps) {
  return (
    <FieldSet label="Display Options">
      <Stack direction="row" gap={1}>
        <InlineField label="Include Group" labelWidth={16}>
          <InlineSwitch
            id={`query-editor-include-group-${query.refId}`}
            value={query.includeGroupName || false}
            onChange={onIncludeGroupName}
          />
        </InlineField>
        <InlineField label="Include Device" labelWidth={16}>
          <InlineSwitch
            id={`query-editor-include-device-${query.refId}`}
            value={query.includeDeviceName || false}
            onChange={onIncludeDeviceName}
          />
        </InlineField>
        <InlineField label="Include Sensor" labelWidth={16}>
          <InlineSwitch
            id={`query-editor-include-sensor-${query.refId}`}
            value={query.includeSensorName || false}
            onChange={onIncludeSensorName}
          />
        </InlineField>
      </Stack>
    </FieldSet>
  );
}
