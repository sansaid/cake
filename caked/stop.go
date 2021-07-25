package main

import (
	"context"
	"fmt"

	pb "github.com/sansaid/cake/pb"
	"github.com/spf13/cobra"
)

// gRPC server methods - this should only get called by the gRPC client (should never be called directly in this code)
func (c *Cake) StopContainer(ctx context.Context, container *pb.Container) (*pb.ContainerStatus, error) {
	fmt.Printf("stopping container: %#v", container)

	return &pb.ContainerStatus{
		Status:      0,
		ContainerId: "ABC123",
		Message:     "Container successfully stopped",
	}, nil
}

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
