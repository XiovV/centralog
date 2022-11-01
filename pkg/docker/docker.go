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

	containerId := c.FindContainerIDByName(container)

	out, _ := c.cli.ContainerLogs(c.ctx, containerId, options)

	stdcopy.StdCopy(logWriter, logWriter, out)
}

func (c *Controller) FindContainerIDByName(name string) string {
	containers, _ := c.cli.ContainerList(c.ctx, types.ContainerListOptions{All: true})

	// [1:] is used because the container name contains a slash at the beginning.
	for _, container := range containers {
		if container.Names[0][1:] == name {
			return container.ID
		}
	}

	return ""
}
