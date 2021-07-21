package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	dockerClient "github.com/docker/docker/client"
	pb "github.com/sansaid/cake/pb"
	"google.golang.org/grpc"
)

const (
	// TODO: confirm the right ports to use - https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml?search=6010
	port = 6010
)

type Cake struct {
	pb.UnimplementedCakedServer
	DockerClient      *dockerClient.Client
	ContainersRunning map[string]int
	StopTimeout       time.Duration
}

// This should only get called by the gRPC client (should never be called directly in this code)
func (c *Cake) StartContainer(ctx context.Context, container *pb.Container) (*pb.ContainerStatus, error) {
	fmt.Printf("starting container: %#v", container)

	return &pb.ContainerStatus{
		Status:      0,
		ContainerId: "ABC123",
		Message:     "Container successfully started",
	}, nil
}

func (c *Cake) StopContainer(ctx context.Context, container *pb.Container) (*pb.ContainerStatus, error) {
	fmt.Printf("stopping container: %#v", container)

	return &pb.ContainerStatus{
		Status:      0,
		ContainerId: "ABC123",
		Message:     "Container successfully stopped",
	}, nil
}

func NewCake() *Cake {
	client, err := dockerClient.NewEnvClient()

	if err != nil {
		panic(err)
	}

	return &Cake{
		DockerClient: client,
		StopTimeout:  30 * time.Second,
	}
}

func main() {
	// TODO: implement start flag for caked
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterCakedServer(grpcServer, NewCake())

	grpcServer.Serve(lis)
}
