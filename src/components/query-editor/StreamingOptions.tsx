import React, { ChangeEvent } from 'react';
import { FieldSet, InlineField, InlineSwitch, Input, Stack } from '@grafana/ui';
import { MyQuery } from '../../types';

interface StreamingOptionsProps {
  query: MyQuery;
  streamIntervalValue: string;
  onStreamingToggle: (event: ChangeEvent<HTMLInputElement>) => void;
  onStreamIntervalChange: (event: ChangeEvent<HTMLInputElement>) => void;
  onStreamIntervalBlur: () => void;
}

export function StreamingOptions({
  query,
  streamIntervalValue,
  onStreamingToggle,
  onStreamIntervalChange,
  onStreamIntervalBlur,
}: StreamingOptionsProps) {
  return (
    <FieldSet label="Streaming Options">
      <Stack direction="row" gap={1}>
        <InlineField label="Enable Streaming" labelWidth={16}>
          <InlineSwitch id="query-editor-is-stream" value={query.isStreaming || false} onChange={onStreamingToggle} />
        </InlineField>
        {query.isStreaming && (
          <InlineField label="Update Interval (ms)" labelWidth={20} tooltip="Refresh interval in milliseconds">
            <Input
              id="query-editor-stream-interval"
              type="number"
              value={streamIntervalValue}
              onChange={onStreamIntervalChange}
              onBlur={onStreamIntervalBlur}
              placeholder="2500"
              min={0}
              max={60000}
            />
          </InlineField>
        )}
      </Stack>
    </FieldSet>
  );
}
