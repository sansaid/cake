package main

import (
	"fmt"
	"net"
	"time"

	dockerClient "github.com/docker/docker/client"
	pb "github.com/sansaid/cake/pb"
	"github.com/sansaid/cake/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var grpcServer *grpc.Server

func NewCake() *Cake {
	client, err := dockerClient.NewEnvClient()

	utils.Check(err, "Cannot create cake client")

	return &Cake{
		DockerClient: client,
		StopTimeout:  30 * time.Second,
	}
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start caked",
	Long:  `Start running the cake daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		defer lis.Close()

		utils.Check(err, "Cannot initialise listener")

		var opts []grpc.ServerOption

		grpcServer = grpc.NewServer(opts...)

		pb.RegisterCakedServer(grpcServer, NewCake())

		grpcServer.Serve(lis)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
