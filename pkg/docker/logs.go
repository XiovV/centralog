package docker

import (
	"fmt"
)

type Log struct {
	ContainerID string
	Message     string
}

func (l *Log) String() string {
	return fmt.Sprintf("container: %s message: %s", l.ContainerID, l.Message)
}

type LogWriter struct {
	buf         *LogBuffer
	containerId string
}

func (l *LogWriter) Write(p []byte) (int, error) {
	str := string(p)
	l.buf.WriteLog(l.containerId, str)

	fmt.Printf("container: %s message: %s", l.containerId, str)

	return len(p), nil
}

func NewLogWriter(buf *LogBuffer, containerId string) *LogWriter {
	return &LogWriter{buf: buf, containerId: containerId}
}
