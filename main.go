package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startImageFlag := startCmd.String("image", "", "[Required] The target image for the watched container, including its register. Example: sansaid/debotbot:latest")

	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	stopImageFlag := stopCmd.String("image", "", "[Required] The target image to stop watching")

	if len(os.Args) < 2 {
		fmt.Errorf("Expected at least one of the following subcommands: start, stop")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		cake := cmds.NewConfig(*image)
		cake.Run()
	case "stop":
		stopCmd.Parse(os.Args[2:])
	default:
		fmt.Errorf("Expected at least one of the following subcommands: start, stop")
		os.Exit(1)
	}
}
