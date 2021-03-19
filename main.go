package main

import (
	"flag"
	"os"

	"github.com/sansaid/cake/cmds"
)

func main() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startImage := startCmd.String("image", "", "[Required] The target image for the watched container, including its register. Example: sansaid/debotbot:latest")

	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)

	if len(os.Args) < 2 {
		cmds.ErrExit(cmds.ErrUnrecognisedSubcommands)
	}

	var cake *Cake

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		cake = NewCake(*startImage)
		cake.Run()
	case "stop":
		stopCmd.Parse(os.Args[2:])

		if cake != nil {
			cake.Stop()
		} else {
			cmds.ErrExit(cmds.ErrCakeNotFound)
		}
	default:
		cmds.ErrExit(cmds.ErrUnrecognisedSubcommands)
	}
}
