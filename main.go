package main

import (
	"flag"
	"os"

	"github.com/sansaid/cake/lib"
)

func main() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startImage := startCmd.String("image", "", "[Required] The target image for the watched container Example: sansaid/debotbot")

	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)

	if len(os.Args) < 2 {
		lib.ExitErr(lib.ErrUnrecognisedSubcommands, nil)
	}

	var cake *lib.Cake

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		cake = lib.NewCake(*startImage, "", "")
		cake.Run()
	case "stop":
		stopCmd.Parse(os.Args[2:])

		if cake != nil {
			cake.Stop()
		} else {
			lib.ExitErr(lib.ErrCakeNotFound, nil)
		}
	default:
		lib.ExitErr(lib.ErrUnrecognisedSubcommands, nil)
	}
}
