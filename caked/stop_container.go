package main

import (
	"context"
	"fmt"

	pb "github.com/sansaid/cake/pb"
	"github.com/sansaid/cake/utils"
	"github.com/spf13/cobra"
)

var stopContainer = func(cake *Cake, id string) {
	ctx := context.TODO()

	err := cake.DockerClient.ContainerStop(ctx, id, &cake.StopTimeout)

	utils.Check(err, "Could not issue stop to container")

	_, errC := cake.DockerClient.ContainerWait(ctx, id, "removed")

	select {
	case err := <-errC:
		utils.Check(err, "Error waiting for container to be removed")
	default:
		return
	}
}

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
