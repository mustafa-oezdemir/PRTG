import { DataQueryRequest, DataQueryResponse, LiveChannelScope } from '@grafana/data';
import { getGrafanaLiveSrv } from '@grafana/runtime';
import { Observable, throwError } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { MyQuery } from '../types';

export function getStreamId(query: MyQuery & { panelId?: string }): string {
  const components = [
    query.panelId || 'default',
    query.refId || 'A',
    query.sensorId || '',
    Array.isArray(query.channelArray) && query.channelArray.length > 0 ? query.channelArray.join('_') : query.channel || '',
  ];

  return components.filter(Boolean).join('_');
}

export function createStreamingQueryObservable(
  datasourceUid: string,
  request: DataQueryRequest<MyQuery>,
  query: MyQuery
): Observable<DataQueryResponse> {
  const queryWithPanelId = {
    ...query,
    panelId: request.panelId?.toString(),
  };

  const streamId = getStreamId(queryWithPanelId);
  const streamPath = `prtg-stream/${streamId}`;

  return getGrafanaLiveSrv()
    .getDataStream({
      addr: {
        scope: LiveChannelScope.DataSource,
        namespace: datasourceUid,
        path: streamPath,
        data: {
          ...query,
          streamId,
          panelId: request.panelId?.toString(),
          queryId: query.refId,
          timeRange: {
            from: request.range.from.valueOf(),
            to: request.range.to.valueOf(),
          },
          cacheTime: query.cacheTime,
          updateMode: query.updateMode,
          bufferSize: query.bufferSize,
        },
      },
    })
    .pipe(
      map((response) => {
        const frameData = response.data || [];
        frameData.forEach((frame) => {
          if (frame && frame.meta) {
            frame.meta = {
              ...frame.meta,
              streaming: true,
              streamId,
              preferredVisualisationType: 'graph',
            };
          }
        });
        return { data: frameData };
      }),
      catchError((err) => {
        console.error('Stream error:', err);
        return throwError(() => new Error(`Streaming error: ${err.message || 'Unknown error'}`));
      })
    );
}
