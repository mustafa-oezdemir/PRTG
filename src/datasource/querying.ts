import { DataQueryRequest, DataQueryResponse } from '@grafana/data';
import { Observable, from, merge, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { MyQuery, QueryType } from '../types';
import { createStreamingQueryObservable } from './streaming';

export type BackendQueryRunner = (request: DataQueryRequest<MyQuery>) => Observable<DataQueryResponse>;

export function filterPrtgQuery(query: MyQuery): boolean {
  if (query.queryType === QueryType.Metrics) {
    return Boolean(query.sensorId && (query.channel || query.channelArray?.length));
  }

  return true;
}

export function runPrtgDataQuery(
  datasourceUid: string,
  request: DataQueryRequest<MyQuery>,
  runBackendQuery: BackendQueryRunner
): Observable<DataQueryResponse> {
  const streamingTargets = request.targets.filter(
    (query) => query.isStreaming && query.queryType === QueryType.Metrics
  );
  const regularTargets = request.targets.filter(
    (query) => !query.isStreaming || query.queryType !== QueryType.Metrics
  );
  const observables: Array<Observable<DataQueryResponse>> = [];

  if (streamingTargets.length > 0) {
    streamingTargets.forEach((query) => {
      observables.push(createStreamingQueryObservable(datasourceUid, request, query));
    });
  }

  if (regularTargets.length > 0) {
    observables.push(
      runBackendQuery({
        ...request,
        targets: regularTargets,
      }).pipe(
        catchError((err) => {
          console.error('Query error:', err);
          return throwError(() => err);
        })
      )
    );
  }

  if (observables.length === 0) {
    return from([{ data: [] }]);
  }

  return merge(...observables);
}
