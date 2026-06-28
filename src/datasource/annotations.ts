import { AnnotationEvent, DataFrame } from '@grafana/data';
import { Observable, from } from 'rxjs';

export function processAnnotationEvents(anno: any, data: DataFrame[]): Observable<AnnotationEvent[]> {
  const events: AnnotationEvent[] = [];

  data.forEach((frame) => {
    const timeField = frame.fields.find((field) => field.name === 'Time');
    const valueField = frame.fields.find((field) => field.name === 'Value');

    if (timeField && valueField) {
      const firstTime = timeField.values[0];
      const lastTime = timeField.values[timeField.values.length - 1];
      const firstValue = valueField.values[0];
      const panelId = typeof anno.panelId === 'number' ? anno.panelId : undefined;
      const source = frame.name || 'PRTG Channel';

      events.push({
        time: firstTime,
        timeEnd: lastTime !== firstTime ? lastTime : undefined,
        title: source,
        text: `Value: ${firstValue}`,
        tags: ['prtg', `value:${firstValue}`, `source:${source}`],
        panelId,
      });
    }
  });

  return from([events]);
}
