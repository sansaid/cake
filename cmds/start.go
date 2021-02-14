package cmds

import (
	"fmt"
)

// NOTE: Use https://github.com/heroku/docker-registry-client

type StartCommand struct {
	*Config
}

// Help - returns the help text for the StartCommand; implements the cli.Command interface
func (s StartCommand) Help() string {
	return "Start a container made by the first image found matching the image name" +
		" and tag identifier. Cake will continue to poll for changes associated" +
		" with this image/tag identifier combo in the registry. If it sees any" +
		" new changes, cake will immediately pull the new image and restart the container."
}

// Run - runs a set of arguments associated with the StartCommand; implements the cli.Command interface
func (s StartCommand) Run(args []string) int {
	runStartCommand(args)

	return 0
}

// Synopsis - synopsis of the StartCommand; implements the cli.Command interface
func (s StartCommand) Synopsis() string {
	return "Start the poll for the specified container image"
}

func runStartCommand(args []string) {
	fmt.Printf("Cake is running with args: %v", args)
}
