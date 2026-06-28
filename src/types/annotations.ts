import { MyQuery } from './query';

export interface AnnotationsQuery extends MyQuery {
  from?: number;
  to?: number;
  limit?: number;
  alertId?: number;
  dashboardId?: number;
  dashboardUID?: string;
  panelId?: number;
  userId?: number;
  type?: 'alert' | 'annotation';
  tags?: string[];
}

export interface Annotation {
  id?: number;
  alertId?: number;
  dashboardId?: number;
  dashboardUID?: string;
  panelId?: number | string;
  userId?: number;
  time: number;
  timeEnd?: number;
  title: string;
  text: string;
  tags?: string[];
  type?: 'alert' | 'annotation';
  data?: Record<string, any>;
}

export interface AnnotationResponse {
  annotations: Annotation[];
  total: number;
}
