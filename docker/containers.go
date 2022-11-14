package docker

import (
	"errors"
	"github.com/docker/docker/api/types"
)

func (c *Controller) GetContainers() ([]types.Container, error) {
	return c.cli.ContainerList(c.ctx, types.ContainerListOptions{All: true})
}

func (c *Controller) GetContainer(name string) (types.Container, error) {
	containers, _ := c.cli.ContainerList(c.ctx, types.ContainerListOptions{All: true})

	for _, container := range containers {
		if container.Names[0][1:] == name {
			return container, nil
		}
	}

	return types.Container{}, errors.New("container could not be found")
}
