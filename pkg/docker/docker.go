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
	logBuffer *LogBuffer
}

func New(db *repository.Repository) *Controller {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &Controller{
		ctx:       ctx,
		cli:       cli,
		logBuffer: NewLogBuffer(db),
	}
}

func (c *Controller) CollectLogs(container string, options types.ContainerLogsOptions) {
	logWriter := NewLogWriter(c.logBuffer, container)

	out, _ := c.cli.ContainerLogs(c.ctx, container, options)

	stdcopy.StdCopy(logWriter, logWriter, out)
}
