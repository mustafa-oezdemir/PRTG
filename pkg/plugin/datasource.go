package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/models"
	healthsvc "github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/health"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtg"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtgtime"
	querysvc "github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/query"
	resourcesvc "github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/resource"
	streamsvc "github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/stream"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
	_ backend.CallResourceHandler   = (*Datasource)(nil)
	_ backend.StreamHandler         = (*Datasource)(nil)
)

func NewDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	config, err := models.LoadPluginSettings(settings)
	if err != nil {
		return nil, err
	}

	prtgtime.SetDefaultTimezone(config.Timezone)

	var cacheTime time.Duration = 60 * time.Second
	if config.CacheTime > 0 {
		cacheTime = config.CacheTime * time.Second
	}

	baseURL := fmt.Sprintf("https://%s", config.Path)
	logger := observability.NewLogger()
	tracer := observability.NewTracer(logger)
	metrics := observability.NewMetrics(prometheus.DefaultRegisterer)
	api := prtg.NewApi(baseURL, config.Secrets.ApiKey, cacheTime, 10*time.Second)
	queryService := querysvc.NewService(api, logger, tracer, metrics, cacheTime)

	ds := &Datasource{
		api:     api,
		logger:  logger,
		tracer:  tracer,
		metrics: metrics,
		query:   queryService,
	}
	ds.resource = resourcesvc.NewService(api, logger, tracer, metrics)
	ds.stream = streamsvc.NewService(logger, queryService, cacheTime)
	ds.health = healthsvc.NewService(api, logger, ds.ClearAllCaches)

	return ds, nil
}

func (d *Datasource) Dispose() {
	if d.query != nil {
		d.query.ClearCache()
	}
	if apiImpl, ok := d.api.(*prtg.Api); ok {
		apiImpl.ClearCache()
	}
}

func (d *Datasource) ClearAllCaches() {
	if d.query != nil {
		d.query.ClearCache()
	}
	if apiImpl, ok := d.api.(*prtg.Api); ok {
		apiImpl.ClearCache()
	}

	d.logger.Debug("All caches cleared")
}
