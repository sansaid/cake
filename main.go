package main

import (
	"factories"

	"github.com/mitchellh/cli"
)

// NOTE: https://pkg.go.dev/github.com/mitchellh/cli#CLI
// NOTE: https://github.com/mitchellh/cli

type cliFactory map[string]cli.CommandFactory

func main() {
	cake := cli.NewCli("cake", "0.1.0")

	cake.Commands = cliFactory{
		"start": factories.StartFactory,
	}

}
