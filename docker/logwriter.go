package docker

import (
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

type BackgroundLogWriter struct {
	buf         *LogBuffer
	containerId string
}

func (l *BackgroundLogWriter) Write(p []byte) (int, error) {
	str := string(p)

	log := parseLog(str, l.containerId)

	l.buf.WriteLog(log)

	//fmt.Printf("containerId: %s message: %s", l.containerId, log.Message)

	return len(p), nil
}

func NewBackgroundLogWriter(buf *LogBuffer, containerId string) *BackgroundLogWriter {
	return &BackgroundLogWriter{buf: buf, containerId: containerId}
}

type ServerLogWriter struct {
	containerId string
	stream      pb.Centralog_FollowLogsServer
}

func (l *ServerLogWriter) Write(p []byte) (int, error) {
	str := string(p)

	log := parseLog(str, l.containerId)

	l.stream.Send(&pb.Log{
		Container: log.ContainerID,
		Timestamp: log.Timestamp,
		Message:   log.Message,
	})

	return len(p), nil
}

func NewServerLogWriter(container string, stream pb.Centralog_FollowLogsServer) *ServerLogWriter {
	return &ServerLogWriter{containerId: container, stream: stream}
}
