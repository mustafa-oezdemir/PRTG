export interface StreamingConfig {
  isStreaming?: boolean;
  streamInterval?: number;
  streamId?: string;
  cacheTime?: number;
  bufferSize?: number;
  updateMode?: 'full' | 'append';
}

export interface StreamingOptions {
  bufferSize: number;
  updateMode: 'full' | 'append';
  cacheTime: number;
  maxStreamsPerPanel: number;
  defaultInterval: number;
}

export const DEFAULT_STREAMING_OPTIONS: StreamingOptions = {
  bufferSize: 100,
  updateMode: 'append',
  cacheTime: 6000,
  maxStreamsPerPanel: 10,
  defaultInterval: 1000,
};
