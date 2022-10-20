package docker

import (
	"fmt"
	"github.com/XiovV/centralog-agent/repository"
)

var (
	StdOut = NewStdout(repository.New())
)

type logMessage map[string]interface{}

type Log struct {
	db  *repository.Repository
	buf *LogBuffer
}

func (l Log) Write(p []byte) (int, error) {
	str := string(p)
	l.buf.AddLog(str)
	fmt.Println(l.buf.Available())

	//fmt.Print("len: ", len(p), " message: ", str)

	return len(p), nil
}

func NewStdout(db *repository.Repository) Log {
	return Log{db: db, buf: NewLogBuffer()}
}
