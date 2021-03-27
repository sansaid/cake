package main

import (
	"flag"
	"os"

	"github.com/sansaid/cake/lib"
)

func main() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startImage := startCmd.String("image", "", "[Required] The target image for the watched container, including its register. Example: sansaid/debotbot:latest")

	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)

	if len(os.Args) < 2 {
		lib.ErrExit(lib.ErrUnrecognisedSubcommands)
	}

	var cake *Cake

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		cake = lib.NewCake(*startImage)
		cake.Run()
	case "stop":
		stopCmd.Parse(os.Args[2:])

		if cake != nil {
			cake.Stop()
		} else {
			lib.ErrExit(lib.ErrCakeNotFound)
		}
	default:
		lib.ErrExit(lib.ErrUnrecognisedSubcommands)
	}
}
