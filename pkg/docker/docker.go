package docker

import (
	"context"
	"github.com/XiovV/centralog-agent/repository"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Controller struct {
	ctx       context.Context
	cli       *client.Client
	logWriter *LogWriter
}

func New() *Controller {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	logWriter := NewLogWriter()

	return &Controller{
		ctx:       ctx,
		cli:       cli,
		logWriter: logWriter,
	}
}

func (c *Controller) GetLogs(container string, options types.ContainerLogsOptions) []repository.LogMessage {
	out, _ := c.cli.ContainerLogs(c.ctx, "d64552a3f96cb6342c06aa055c5579046b6da94f362ef2360ad5aefadc3d05b1", options)

	stdcopy.StdCopy(c.logWriter, c.logWriter, out)
	return c.logWriter.GetLogs()
}
