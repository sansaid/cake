package cmds

import (
	"flag"
)

func setupStopCmd() {
	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	stopCmd.String("image", "", "[Required] The target image to stop watching")
}
