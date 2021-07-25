package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	dockerClient "github.com/docker/docker/client"
	pb "github.com/sansaid/cake/pb"
	"github.com/sansaid/cake/utils"
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

// gRPC server methods - this should only get called by the gRPC client (should never be called directly in this code)
func (c *Cake) StartContainer(ctx context.Context, container *pb.Container) (*pb.ContainerStatus, error) {
	fmt.Printf("starting container: %#v", container)

	return &pb.ContainerStatus{
		Status:      0,
		ContainerId: "ABC123",
		Message:     "Container successfully started",
	}, nil
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

func NewCake() *Cake {
	client, err := dockerClient.NewEnvClient()

	utils.Check(err, "Cannot create cake client")

	return &Cake{
		DockerClient: client,
		StopTimeout:  30 * time.Second,
	}
}

func main() {
	// TODO: implement start flag for caked
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	defer lis.Close()

	utils.Check(err, "Cannot initialise listener")

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterCakedServer(grpcServer, NewCake())

	grpcServer.Serve(lis)
}
