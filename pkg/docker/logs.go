package docker

import (
	"fmt"
	"strings"
	"time"
)

type Log struct {
	ContainerID string
	Message     string
}

type LogWriter struct {
	buf         *LogBuffer
	containerId string
}

func (l *LogWriter) Write(p []byte) (int, error) {
	str := string(p)

	strParts := strings.SplitN(str, " ", 2)

	log := strParts[1]
	timestamp, err := time.Parse("2006-01-02T15:04:05Z", strParts[0])
	if err != nil {
		fmt.Println(err)
	}

	l.buf.WriteLog(l.containerId, log, timestamp.UnixMilli())

	fmt.Printf("container: %s message: %s", l.containerId, log)

	return len(p), nil
}

func NewLogWriter(buf *LogBuffer, containerId string) *LogWriter {
	return &LogWriter{buf: buf, containerId: containerId}
}
