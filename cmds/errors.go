package cmds

import "fmt"

// ErrNoRegistry - an error used to specify no if no registry is defined when running cake
var ErrNoRegistry error = fmt.Errorf("Registry must be specified, either through -registry flag, CAKE_REGISTRY" +
	" environment variable, or through a .cake file.")

// ErrUnrecognisedSubcommands - an error to specify subcommands are uncregonised
var ErrUnrecognisedSubcommands error = fmt.Errorf("Expected at least one of the following subcommands: start, stop")
