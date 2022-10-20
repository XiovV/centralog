package docker

import (
	"fmt"
	"github.com/XiovV/centralog-agent/repository"
)

var (
	StdOut = NewStdout(repository.New())
)

type Log struct {
	db *repository.Repository
}

func (l Log) Write(p []byte) (int, error) {
	str := string(p)

	fmt.Print("message: ", str)

	return len(p), nil
}

func NewStdout(db *repository.Repository) Log {
	return Log{db: db}
}
