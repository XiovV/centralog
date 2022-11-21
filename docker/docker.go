package docker

import (
	"context"
	"fmt"
	"github.com/XiovV/centralog-agent/repository"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
)

type Controller struct {
	ctx       context.Context
	cli       *client.Client
	logBuffer *LogBuffer
}

func New(db *repository.SQLite) *Controller {
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

func (c *Controller) CollectLogs(container string, writer io.Writer, options types.ContainerLogsOptions) {
	containerId := c.FindContainerIDByName(container)

	out, _ := c.cli.ContainerLogs(c.ctx, containerId, options)

	stdcopy.StdCopy(writer, writer, out)
}

func (c *Controller) CollectLogsBackground(container string, writer io.Writer, options types.ContainerLogsOptions) {
	containerId := c.FindContainerIDByName(container)

	out, _ := c.cli.ContainerLogs(c.ctx, containerId, options)

	_, err := stdcopy.StdCopy(writer, writer, out)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (c *Controller) FindContainerIDByName(name string) string {
	containers, _ := c.cli.ContainerList(c.ctx, types.ContainerListOptions{All: true})

	// [1:] is used because the containerId name contains a slash at the beginning.
	for _, container := range containers {
		if container.Names[0][1:] == name {
			return container.ID
		}
	}

	return ""
}
