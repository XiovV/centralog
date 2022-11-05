package docker

import (
	"github.com/XiovV/centralog-agent/repository"
)

const (
	defaultBufferSize = 512
)

type LogBuffer struct {
	buf []repository.LogMessage
	db  *repository.Repository
}

func NewLogBuffer(db *repository.Repository) *LogBuffer {
	return &LogBuffer{buf: []repository.LogMessage{}, db: db}
}

func (b *LogBuffer) WriteLog(log Log) {
	if len(b.buf) == defaultBufferSize-1 {
		b.buf = append(b.buf, repository.LogMessage{"containerID": log.ContainerID, "message": log.Message, "timestamp": log.Timestamp})
		b.db.StoreLogs(b.buf)
		b.buf = nil
	} else {
		b.buf = append(b.buf, repository.LogMessage{"containerID": log.ContainerID, "message": log.Message, "timestamp": log.Timestamp})
	}
}

func (b *LogBuffer) Flush() {
	if len(b.buf) > 0 {
		b.db.StoreLogs(b.buf)
		b.buf = nil
	}
}
