package stream

import "time"

func updateChannelBuffer(stream *activeStream, state *channelState, times []time.Time, values []float64) {
	if len(values) > 0 {
		state.lastValue = values[len(values)-1]
	}

	if stream.updateMode == "append" {
		if len(times) > 0 {
			appendBufferData(state.buffer, times, values, stream.bufferSize)
		}
	} else {
		state.buffer.times = times
		state.buffer.values = values

		if int64(len(state.buffer.times)) > state.buffer.size {
			excess := int64(len(state.buffer.times)) - state.buffer.size
			state.buffer.times = state.buffer.times[excess:]
			state.buffer.values = state.buffer.values[excess:]
		}
	}
}

func appendBufferData(buffer *dataBuffer, newTimes []time.Time, newValues []float64, maxSize int64) {
	curLen := len(buffer.times)
	newLen := curLen + len(newTimes)

	if int64(len(newTimes)) >= maxSize {
		startIdx := len(newTimes) - int(maxSize)
		buffer.times = newTimes[startIdx:]
		buffer.values = newValues[startIdx:]
		return
	}

	if int64(newLen) > maxSize {
		excess := int64(newLen) - maxSize
		buffer.times = buffer.times[int(excess):]
		buffer.values = buffer.values[int(excess):]
	}

	buffer.times = append(buffer.times, newTimes...)
	buffer.values = append(buffer.values, newValues...)
}
