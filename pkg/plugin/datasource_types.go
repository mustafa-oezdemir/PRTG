package plugin

import (
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/health"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtg"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/query"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/resource"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/stream"
)

type Datasource struct {
	api      prtg.PRTGAPI
	logger   observability.PrtgLogger
	tracer   *observability.Tracer
	metrics  *observability.Metrics
	query    *query.Service
	resource *resource.Service
	stream   *stream.Service
	health   *health.Service
}
