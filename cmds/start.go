package cmds

import (
	"flag"
)

func setupStartCmd() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	imageFlag := startCmd.String("image", "", "[Required] The target image for the watched container, including its register. Example: sansaid/debotbot:latest")
}
