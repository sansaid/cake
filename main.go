package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/sansaid/cake/cmds"
)

// NOTE: https://pkg.go.dev/github.com/mitchellh/cli#CLI
// NOTE: https://github.com/mitchellh/cli

type cliFactory map[string]cli.CommandFactory

func main() {
	cake := cli.NewCLI("cake", "0.1.0")

	cake.Args = os.Args[1:]
	config, _ := cmds.ParseArgs(cake.Args)

	cake.Commands = cliFactory{
		"start": func() (cli.Command, error) {
			return cmds.StartCommand{
				Config: config,
			}, nil
		},
		"stop": func() (cli.Command, error) {
			return cmds.StopCommand{
				Config: config,
			}, nil
		},
	}

	exitStatus, err := cake.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
