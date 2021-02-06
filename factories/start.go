package factories

import (
	"github.com/mitchellh/cli"
)

type ContainerRef struct {
	repository string
	tag        string
}

// NOTE: Use https://github.com/heroku/docker-registry-client

//
func (c *ContainerRef) Help() string {}

//
func (c *ContainerRef) Run(args []string) int {}

//
func (c *ContainerRef) Synopsis() string {
	return "Start the poll for the specified container image"
}

// NewContainerRef - creates a new ContainerRef struct
func NewContainerRef(repository string, tag string) *ContainerRef {
	return &ContainerRef{
		repository: repository,
		tag:        tag,
	}
}

// StartFactory - a CommandFactory implementation, which returns an implementation of the cli.Command interface
func StartFactory() (cli.Command, error) {
	return
}
