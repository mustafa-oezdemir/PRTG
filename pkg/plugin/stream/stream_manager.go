package stream

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Service) getExistingStream(streamID string) *activeStream {
	s.streamManager.mu.RLock()
	existingStream, exists := s.streamManager.streams[streamID]
	s.streamManager.mu.RUnlock()

	if !exists {
		return nil
	}
	return existingStream
}

func (s *Service) updateExistingStream(stream *activeStream, from time.Time, to time.Time, cacheDuration time.Duration) {
	s.streamManager.mu.Lock()
	stream.timeRange = &backend.TimeRange{From: from, To: to}
	stream.isActive = true
	stream.lastUpdate = time.Now().Add(-cacheDuration)
	s.streamManager.mu.Unlock()

	select {
	case stream.updateChan <- struct{}{}:
	default:
	}
}

func (s *Service) registerNewStream(stream *activeStream, streamID string) {
	s.streamManager.mu.Lock()
	s.streamManager.streams[streamID] = stream
	s.streamManager.mu.Unlock()

	s.trackStream(stream.panelId, streamID, stream)
}

func (s *Service) cleanupStream(stream *activeStream) {
	s.streamManager.mu.Lock()
	if currentStream, exists := s.streamManager.streams[stream.streamID]; exists && currentStream == stream {
		delete(s.streamManager.streams, stream.streamID)
	}

	if panelStreams, exists := s.streamManager.activeStreams[stream.panelId]; exists {
		if currentStream, exists := panelStreams[stream.streamID]; exists && currentStream == stream {
			delete(panelStreams, stream.streamID)
		}
		if len(panelStreams) == 0 {
			delete(s.streamManager.activeStreams, stream.panelId)
		}
	}
	s.streamManager.mu.Unlock()

	s.logger.Debug("Stream closed", "streamID", stream.streamID)
}

func (s *Service) trackStream(panelId string, streamId string, stream *activeStream) {
	s.streamManager.mu.Lock()
	defer s.streamManager.mu.Unlock()

	if _, exists := s.streamManager.activeStreams[panelId]; !exists {
		s.streamManager.activeStreams[panelId] = make(map[string]*activeStream)
	}

	s.streamManager.streams[streamId] = stream
	s.streamManager.activeStreams[panelId][streamId] = stream

	s.logger.Debug("Stream tracked",
		"panelId", panelId,
		"streamId", streamId,
		"totalStreams", len(s.streamManager.streams),
		"panelStreams", len(s.streamManager.activeStreams[panelId]))
}

func (s *Service) getStreamsByPanel(panelId string) []*activeStream {
	s.streamManager.mu.RLock()
	defer s.streamManager.mu.RUnlock()

	result := make([]*activeStream, 0, 5)
	if panelStreams, exists := s.streamManager.activeStreams[panelId]; exists {
		for _, stream := range panelStreams {
			result = append(result, stream)
		}
	}
	return result
}
