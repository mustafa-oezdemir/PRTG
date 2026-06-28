import React from 'react';
import { Combobox, FieldSet, InlineField, Stack, type ComboboxOption } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

interface PropertyOptionsProps {
  property?: string;
  filterProperty?: string;
  properties: Array<SelectableValue<string>>;
  filterProperties: Array<SelectableValue<string>>;
  onPropertyChange: (value: string) => void;
  onFilterPropertyChange: (value: string) => void;
}

export function PropertyOptions({
  property,
  filterProperty,
  properties,
  filterProperties,
  onPropertyChange,
  onFilterPropertyChange,
}: PropertyOptionsProps) {
  return (
    <FieldSet label="Options">
      <Stack direction="row" gap={2}>
        <InlineField label="Property" labelWidth={16} tooltip="Select property type">
          <Combobox
            id="query-editor-property"
            options={properties.map((p) => ({ label: p.label!, value: p.value! }))}
            value={property}
            onChange={(option: ComboboxOption<string> | null) => {
              if (option?.value) {
                onPropertyChange(option.value);
              }
            }}
            width={32}
            placeholder="Select property"
            isClearable={false}
          />
        </InlineField>
        <InlineField label="Filter Property" labelWidth={16} tooltip="Select filter property">
          <Combobox
            id="query-editor-filterProperty"
            options={filterProperties.map((p) => ({ label: p.label!, value: p.value! }))}
            value={filterProperty}
            onChange={(option: ComboboxOption<string> | null) => {
              if (option?.value) {
                onFilterPropertyChange(option.value);
              }
            }}
            width={32}
            placeholder="Select filter"
            isClearable={false}
          />
        </InlineField>
      </Stack>
    </FieldSet>
  );
}
