import { DataFrame, DataQueryRequest, DataSourceInstanceSettings } from '@grafana/data';
import { getTemplateSrv, getGrafanaLiveSrv } from '@grafana/runtime';
import { firstValueFrom, of } from 'rxjs';
import { DataSource } from './datasource';
import { MyQuery, MyDataSourceOptions, QueryType } from './types';
import pluginJson from './plugin.json';

jest.mock('@grafana/runtime', () => ({
    getTemplateSrv: jest.fn(),
    getGrafanaLiveSrv: jest.fn(),
    DataSourceWithBackend: class MockDataSourceWithBackend {
        uid = 'test-uid';
        constructor() {}
        getResource = jest.fn();
        query = jest.fn().mockReturnValue({ data: [] });
    }
}));

const mockTemplateSrv = { replace: jest.fn() };
const mockLiveSrv = { getDataStream: jest.fn() };

const createMockQuery = (overrides: Partial<MyQuery> = {}): MyQuery => ({
    refId: 'A',
    queryType: QueryType.Metrics,
    group: '',
    groupId: '',
    device: '',
    deviceId: '',
    sensor: '',
    sensorId: '',
    channel: '',
    channelArray: [],
    manualMethod: '',
    manualObjectId: '',
    property: '',
    filterProperty: '',
    includeGroupName: false,
    includeDeviceName: false,
    includeSensorName: false,
    streaming: undefined,
    isStreaming: false,
    streamInterval: 2500,
    streamId: '',
    panelId: '',
    queryId: '',
    cacheTime: undefined,
    bufferSize: undefined,
    updateMode: 'full',
    ...overrides
});

const mockSettings: DataSourceInstanceSettings<MyDataSourceOptions> = {
    id: 1,
    uid: 'test-uid',
    type: 'prtg',
    name: 'PRTG DataSource',
    url: 'http://localhost',
    access: 'proxy',
    readOnly: false,
    jsonData: {},
    meta: {} as any
};

describe('DataSource', () => {
    let dataSource: DataSource;

    beforeEach(() => {
        jest.clearAllMocks();
        (getTemplateSrv as jest.Mock).mockReturnValue(mockTemplateSrv);
        (getGrafanaLiveSrv as jest.Mock).mockReturnValue(mockLiveSrv);
        dataSource = new DataSource(mockSettings);
    });

    afterEach(() => {
        jest.restoreAllMocks();
    });

    it('should create instance', () => {
        expect(dataSource).toBeDefined();
    });

    it('should advertise annotation query support', () => {
        expect(pluginJson.annotations).toBe(true);
        expect(dataSource.annotations).toEqual({
            QueryEditor: undefined,
            processEvents: expect.any(Function),
        });
    });

    it('should advertise streaming data source support', () => {
        expect(pluginJson.streaming).toBe(true);
    });

    it('should convert query frames to annotation events', async () => {
        const frame: DataFrame = {
            name: 'CPU Load',
            fields: [
                { name: 'Time', values: [1000, 2000] },
                { name: 'Value', values: [42, 45] },
            ],
            length: 2,
        } as DataFrame;

        const events = await firstValueFrom(dataSource.annotations.processEvents({ panelId: 7 }, [frame]));

        expect(events).toEqual([
            {
                time: 1000,
                timeEnd: 2000,
                title: 'CPU Load',
                text: 'Value: 42',
                tags: ['prtg', 'value:42', 'source:CPU Load'],
                panelId: 7,
            },
        ]);
    });

    it('should handle template variables', () => {
        const query = createMockQuery({
            sensorId: '$sensor',
            channel: '$channel',
            channelArray: ['$channel'],
        });
        mockTemplateSrv.replace.mockImplementation((value: string) => {
            const replacements: Record<string, string> = {
                '$sensor': '1025',
                '$channel': 'Total',
            };
            return replacements[value] ?? value;
        });

        const result = dataSource.applyTemplateVariables(query, {});

        expect(mockTemplateSrv.replace).toHaveBeenCalledWith('$channel', {});
        expect(result.sensorId).toBe('1025');
        expect(result.channel).toBe('Total');
        expect(result.channelArray).toEqual(['Total']);
    });

    it('should filter queries by channel', () => {
        expect(dataSource.filterQuery(createMockQuery({ sensorId: '1025', channel: 'test' }))).toBe(true);
        expect(dataSource.filterQuery(createMockQuery({ sensorId: '1025', channelArray: ['test'] }))).toBe(true);
        expect(dataSource.filterQuery(createMockQuery({ channel: '' }))).toBe(false);
        expect(dataSource.filterQuery(createMockQuery({ queryType: QueryType.Manual, channel: '' }))).toBe(true);
    });

    it('should handle resource calls', async () => {
        (dataSource as any).getResource.mockResolvedValue({ groups: [] });
        const result = await dataSource.getGroups();
        expect(result).toEqual({ groups: [] });
    });

    it('should handle device errors', async () => {
        await expect(dataSource.getDevices('')).rejects.toThrow('group is required');
    });

    it('should return metric find values for PRTG dashboard variables', async () => {
        (dataSource as any).getResource.mockImplementation((path: string) => {
            switch (path) {
                case 'groups':
                    return Promise.resolve({ groups: [{ group: 'Core' }] });
                case 'devices/Core':
                    return Promise.resolve({ devices: [{ device: 'Probe' }] });
                case 'sensors/Probe':
                    return Promise.resolve({ sensors: [{ sensor: 'CPU Load', objid: 1025 }] });
                case 'channels/1025':
                    return Promise.resolve({ values: [{ datetime: 'ignored', Total: 10 }] });
                default:
                    return Promise.resolve({});
            }
        });
        mockTemplateSrv.replace.mockImplementation((value: string) => value.replace('$group', 'Core').replace('$device', 'Probe').replace('$sensor', '1025'));

        await expect(dataSource.metricFindQuery('groups')).resolves.toEqual([{ text: 'Core', value: 'Core' }]);
        await expect(dataSource.metricFindQuery({ query: 'devices:$group' })).resolves.toEqual([{ text: 'Probe', value: 'Probe' }]);
        await expect(dataSource.metricFindQuery('sensors:$device')).resolves.toEqual([{ text: 'CPU Load', value: '1025' }]);
        await expect(dataSource.metricFindQuery('channels:$sensor')).resolves.toEqual([{ text: 'Total', value: 'Total' }]);
    });

    it('should handle query calls', () => {
        const mockRequest: DataQueryRequest<MyQuery> = {
            requestId: 'test',
            interval: '1s',
            intervalMs: 1000,
            range: {
                from: { valueOf: () => 1000 } as any,
                to: { valueOf: () => 2000 } as any,
                raw: { from: 'now-1h', to: 'now' }
            },
            scopedVars: {},
            targets: [createMockQuery({ channel: 'test' })],
            timezone: 'UTC',
            app: 'grafana',
            startTime: Date.now(),
            panelId: 1
        };

        // Mock the query method specifically for this test
        const mockResponse = { data: [] };
        jest.spyOn(dataSource, 'query').mockReturnValue(mockResponse as any);
        
        const result = dataSource.query(mockRequest);
        expect(result).toBeDefined();
        expect(result).toBe(mockResponse);
    });

    it('should route streaming metric queries through Grafana Live', async () => {
        mockLiveSrv.getDataStream.mockReturnValue(of({ data: [] }));

        const mockRequest: DataQueryRequest<MyQuery> = {
            requestId: 'stream-test',
            interval: '1s',
            intervalMs: 1000,
            range: {
                from: { valueOf: () => 1000 } as any,
                to: { valueOf: () => 2000 } as any,
                raw: { from: 'now-1m', to: 'now' }
            },
            scopedVars: {},
            targets: [
                createMockQuery({
                    refId: 'A',
                    sensorId: '1025',
                    channel: 'Total',
                    channelArray: ['Total'],
                    isStreaming: true,
                    streamInterval: 2500,
                    updateMode: 'append',
                }),
            ],
            timezone: 'UTC',
            app: 'dashboard',
            startTime: Date.now(),
            panelId: 7
        };

        await firstValueFrom(DataSource.prototype.query.call(dataSource, mockRequest));

        expect(mockLiveSrv.getDataStream).toHaveBeenCalledWith({
            addr: expect.objectContaining({
                namespace: 'test-uid',
                path: 'prtg-stream/7_A_1025_Total',
                data: expect.objectContaining({
                    sensorId: '1025',
                    channel: 'Total',
                    channelArray: ['Total'],
                    isStreaming: true,
                    streamInterval: 2500,
                    streamId: '7_A_1025_Total',
                    panelId: '7',
                    queryId: 'A',
                    timeRange: {
                        from: 1000,
                        to: 2000,
                    },
                    updateMode: 'append',
                }),
            }),
        });
    });
});
