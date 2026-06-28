package resource

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/observability"
	"github.com/1DeliDolu/PRTG/maxmarkusprogram-prtg-datasource/pkg/plugin/prtg"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

const MaxQueueSize = 100

type Service struct {
	api          prtg.PRTGAPI
	logger       observability.PrtgLogger
	tracer       *observability.Tracer
	metrics      *observability.Metrics
	requestQueue []*ResourceRequest
	queueLock    sync.Mutex
}

func NewService(
	api prtg.PRTGAPI,
	logger observability.PrtgLogger,
	tracer *observability.Tracer,
	metrics *observability.Metrics,
) *Service {
	return &Service{
		api:          api,
		logger:       logger,
		tracer:       tracer,
		metrics:      metrics,
		requestQueue: make([]*ResourceRequest, 0),
	}
}

type ResourceRequest struct {
	Request *backend.CallResourceRequest
	Sender  backend.CallResourceResponseSender
}

func (s *Service) processQueuedRequests() error {
	s.queueLock.Lock()
	defer s.queueLock.Unlock()

	if len(s.requestQueue) == 0 {
		return nil
	}

	if len(s.requestQueue) > MaxQueueSize {
		s.logger.Warn("Request queue overflow, dropping old requests")
		s.requestQueue = s.requestQueue[len(s.requestQueue)-MaxQueueSize:]
	}

	var lastError error
	for _, req := range s.requestQueue {
		err := s.processRequest(req)
		if err != nil {
			s.logger.Error("Failed to process request", "error", err)
			lastError = err
		}
	}

	s.requestQueue = s.requestQueue[:0]
	return lastError
}

func (s *Service) processRequest(req *ResourceRequest) error {
	path := req.Request.Path
	s.logger.Debug("Processing request", "path", path)

	switch {
	case strings.HasPrefix(path, "groups"):
		return s.handleGetGroups(req.Sender)

	case strings.HasPrefix(path, "devices/"):
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			return sendErrorResponse(req.Sender, "group parameter is required", http.StatusBadRequest)
		}
		return s.handleGetDevices(req.Sender, pathParts[1])

	case strings.HasPrefix(path, "sensors/"):
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			return sendErrorResponse(req.Sender, "device parameter is required", http.StatusBadRequest)
		}
		return s.handleGetSensors(req.Sender, pathParts[1])

	case strings.HasPrefix(path, "channels/"):
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			return sendErrorResponse(req.Sender, "sensor parameter is required", http.StatusBadRequest)
		}
		return s.handleGetChannel(req.Sender, pathParts[1])

	default:
		return sendErrorResponse(req.Sender, "invalid API endpoint", http.StatusNotFound)
	}
}

func sendErrorResponse(sender backend.CallResourceResponseSender, message string, statusCode int) error {
	errorResponse := map[string]string{"error": message}
	errorJSON, _ := json.Marshal(errorResponse)
	return sender.Send(&backend.CallResourceResponse{
		Status:  statusCode,
		Headers: map[string][]string{"Content-Type": {"application/json"}},
		Body:    errorJSON,
	})
}
