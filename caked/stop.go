package main

import (
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop caked",
	Long:  `Stop the cake daemon gracefully`,
	Run: func(cmd *cobra.Command, args []string) {
		grpcServer.GracefulStop()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
