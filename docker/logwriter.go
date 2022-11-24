package docker

import (
	"errors"
	pb "github.com/XiovV/centralog-agent/grpc"
	"strings"
	"time"
)

type Log struct {
	ContainerID string
	Message     string
	Timestamp   int64
}

func parseLog(log, containerId string) Log {
	logParts := strings.SplitN(log, " ", 2)

	logMessage := logParts[1]
	timestamp, _ := time.Parse("2006-01-02T15:04:05Z", logParts[0])

	return Log{
		ContainerID: containerId,
		Message:     logMessage,
		Timestamp:   timestamp.UnixMilli(),
	}
}

func NewBackgroundLogWriter(buf *LogBuffer, containerId string) *BackgroundLogWriter {
	return &BackgroundLogWriter{buf: buf, containerId: containerId}
}

type BackgroundLogWriter struct {
	buf         *LogBuffer
	containerId string
}

func (l *BackgroundLogWriter) Write(p []byte) (int, error) {
	str := string(p)

	log := parseLog(str, l.containerId)

	l.buf.WriteLog(log)

	return len(p), nil
}

type ServerLogWriter struct {
	containerId string
	stream      pb.Centralog_GetLogsServer
	stopSignal  chan struct{}
}

func (l *ServerLogWriter) Write(p []byte) (int, error) {
	str := string(p)

	log := parseLog(str, l.containerId)

	l.stream.Send(&pb.Log{
		Container: log.ContainerID,
		Timestamp: log.Timestamp,
		Message:   log.Message,
	})

	select {
	case <-l.stopSignal:
		return 0, errors.New("stopping " + log.ContainerID)
	default:
		return len(p), nil
	}
}

func NewServerLogWriter(stopSignal chan struct{}, container string, stream pb.Centralog_GetLogsServer) *ServerLogWriter {
	return &ServerLogWriter{containerId: container, stream: stream, stopSignal: stopSignal}
}
