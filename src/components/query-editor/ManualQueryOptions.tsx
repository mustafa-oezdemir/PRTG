import React, { ChangeEvent } from 'react';
import { Combobox, FieldSet, InlineField, Input, Stack, type ComboboxOption } from '@grafana/ui';
import { ManualApiMethod } from '../../types';

interface ManualQueryOptionsProps {
  manualMethods: ManualApiMethod[];
  manualMethod: string;
  manualObjectId: string;
  sensorId: string;
  onManualMethodChange: (value: string) => void;
  onManualObjectIdChange: (event: ChangeEvent<HTMLInputElement>) => void;
}

export function ManualQueryOptions({
  manualMethods,
  manualMethod,
  manualObjectId,
  sensorId,
  onManualMethodChange,
  onManualObjectIdChange,
}: ManualQueryOptionsProps) {
  return (
    <FieldSet label="Manual API Query">
      <Stack direction="row" gap={2}>
        <InlineField label="API Method" labelWidth={16} tooltip="Select or enter a custom PRTG API endpoint">
          <Combobox
            id="query-editor-manualMethod"
            options={manualMethods.map((method) => ({
              label: method.label!,
              value: method.value!,
            }))}
            value={manualMethod}
            onChange={(option: ComboboxOption<string> | null) => {
              if (option?.value) {
                onManualMethodChange(option.value);
              }
            }}
            width={32}
            placeholder="Select or enter API method"
            createCustomValue={true}
            isClearable={true}
          />
        </InlineField>
        <InlineField label="Object ID" labelWidth={16} tooltip="Object ID from selected sensor">
          <Input
            id="query-editor-manualObjectId"
            value={manualObjectId || sensorId}
            onChange={onManualObjectIdChange}
            placeholder="Automatically filled from sensor"
            width={32}
            type="text"
            disabled={!!sensorId}
          />
        </InlineField>
      </Stack>
    </FieldSet>
  );
}
