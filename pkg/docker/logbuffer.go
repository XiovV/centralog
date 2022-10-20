package docker

import "fmt"

const (
	defaultLogBufferSize = 10
)

type LogBuffer struct {
	buf []map[string]interface{}
}

func NewLogBuffer() *LogBuffer {
	return &LogBuffer{buf: []map[string]interface{}{}}
}

func (b *LogBuffer) AddLog(message string) {
	if len(b.buf) == defaultLogBufferSize-1 {
		b.buf = append(b.buf, logMessage{"message": message})
		fmt.Println(b.buf) // replace with a sql insert
		b.buf = nil
	} else {
		b.buf = append(b.buf, logMessage{"message": message})
	}
}

func (b *LogBuffer) Flush() {
	if len(b.buf) > 0 {
		fmt.Println(b.buf) // replace with a sql insert
		b.buf = nil
	}
}

func (b *LogBuffer) Available() int {
	return defaultLogBufferSize - len(b.buf)
}

func (b *LogBuffer) Buffered() int {
	return len(b.buf)
}
