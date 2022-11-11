package docker

import "github.com/docker/docker/api/types"

func (c *Controller) GetContainers() ([]types.Container, error) {
	return c.cli.ContainerList(c.ctx, types.ContainerListOptions{All: true})
}
