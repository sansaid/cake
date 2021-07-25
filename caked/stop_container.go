package main

import (
	"context"

	"github.com/sansaid/cake/utils"
)

// TODO: The use of the word `digest` in variable and function names is inconsistently being used - need to make this more consistent
// TODO: Write functionality to sync `cake` with the local system that's being managed externally to it
// TODO: Think about how to deal with pruning containers and images with `cake` on RasbPi - should stopped containers also be deleted and have their images removed?

var closeClient = func(c *Cake) {
	c.DockerClient.Close()
}

var stopContainer = func(cake *Cake, id string) {
	ctx := context.TODO()

	err := cake.DockerClient.ContainerStop(ctx, id, &cake.StopTimeout)

	utils.Check(err, "Could not issue stop to container")

	_, errC := cake.DockerClient.ContainerWait(ctx, id, "removed")

	select {
	case err := <-errC:
		utils.Check(err, "Error waiting for container to be removed")
	default:
		return
	}
}

// Stop - stop this instance of cake and remove all managed containers
func (c *Cake) Stop() {
	defer closeClient(c)

	c.ContainersRunning.Lock()
	for id, _ := range c.ContainersRunning.containers {
		stopContainer(c, id)
		_, exists := c.ContainersRunning.containers[id]

		if exists {
			delete(c.ContainersRunning.containers, id)
		}
	}
	c.ContainersRunning.Unlock()
}
