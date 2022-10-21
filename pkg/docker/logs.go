package docker

import (
	"github.com/XiovV/centralog-agent/repository"
)

type LogWriter struct {
	buf []repository.LogMessage
}

func (l *LogWriter) Write(p []byte) (int, error) {
	str := string(p)
	l.addLog(str)

	return len(p), nil
}

func (l *LogWriter) addLog(message string) {
	l.buf = append(l.buf, repository.LogMessage{"message": message})
}

func (l *LogWriter) GetLogs() []repository.LogMessage {
	return l.buf
}

func NewLogWriter() *LogWriter {
	return &LogWriter{buf: []repository.LogMessage{}}
}
