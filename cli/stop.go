package main

import (
	"context"

	pb "github.com/sansaid/cake/pb"
	"github.com/sansaid/cake/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a cake image",
	Long:  `Stop containers associated with this image. The container will no longer run and you will no longer receive updates from Docker Hub.`,
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption

		conn, err := grpc.Dial("localhost:6010", opts...)
		defer conn.Close()

		utils.Check(err, "Cannot initialise gRPC dial")

		container := NewContainer(Image)
		client := pb.NewCakedClient(conn)

		client.StopContainer(context.Background(), container)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
