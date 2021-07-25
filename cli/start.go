package main

import (
	"context"
	"log"
	"strings"

	pb "github.com/sansaid/cake/pb"
	"github.com/sansaid/cake/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NewContainer(image string) *pb.Container {
	registry := "https://hub.docker.com"

	splitImage := strings.Split(Image, ":")

	if len(splitImage) != 2 {
		log.Fatalf("Image must be of the foramt <repo>/<image>:<tag>")
	}

	repo, tag := splitImage[0], splitImage[1]

	return &pb.Container{
		ImageName: repo,
		Tag:       tag,
		Registry:  registry,
	}
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a cake image",
	Long:  `Start running an image managed by cake. Any changes to the image in Docker Hub will automatically get updated where the corresponding image is running.`,
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption

		conn, err := grpc.Dial("localhost:6010", opts...)
		defer conn.Close()

		utils.Check(err, "Cannot initialise gRPC dial")

		container := NewContainer(Image)
		client := pb.NewCakedClient(conn)

		client.StartContainer(context.Background(), container)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
