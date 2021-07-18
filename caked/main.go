package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/sansaid/cake/pb"
	"google.golang.org/grpc"
)

const (
	// TODO: confirm the right ports to use - https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml?search=6010
	port = 6010
)

type cakedServer struct {
	pb.UnimplementedCakedServer
}

// This should only get called by the gRPC client (should never be called directly in this code)
// TODO: ^^confirm the above is correct
func (c *cakedServer) StartContainer(ctx context.Context, container *pb.Container) (*pb.Started, error) {
	fmt.Printf("starting container: %#v", container)

	return &pb.Started{Started: true}, nil
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

	pb.RegisterCakedServer(grpcServer, &cakedServer{})

	grpcServer.Serve(lis)
}
