package cmds

import "fmt"

type StopCommand struct {
	*Config
}

// Help - returns the help text for the StopCommand; implements the cli.Command interfaces
func (s StopCommand) Help() string {
	return "Stops the container associated with the image and tag identifier, including polling any updates" +
		" from the registry. Use this command if you want to stop the container managed by cake including" +
		" receiving any updates to it."
}

// Run - runs a set of arguments associated with the StopCommand; implements the cli.Command interface
func (s StopCommand) Run(args []string) int {
	runStopCommand(args)

	return 0
}

// Synopsis - a synopsis of the StopCommand; implements the cli.Command interface
func (s StopCommand) Synopsis() string {
	return "Stop the poll for the specified container image - this does not stop the container image running"
}

func runStopCommand(args []string) {
	fmt.Println("Stopping")
}
