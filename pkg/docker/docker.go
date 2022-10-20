package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
)

type Controller struct {
	ctx context.Context
	cli *client.Client
}

func New() *Controller {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &Controller{ctx: ctx, cli: cli}
}

func (c *Controller) GetLogs(container string, options types.ContainerLogsOptions) io.ReadCloser {
	out, _ := c.cli.ContainerLogs(c.ctx, "d64552a3f96cb6342c06aa055c5579046b6da94f362ef2360ad5aefadc3d05b1", options)

	return out
}
