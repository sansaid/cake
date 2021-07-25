package utils

import (
	"log"

	"google.golang.org/grpc"
)

// TODO: delete if not in use
type Runner func(connection *grpc.ClientConn)

func Check(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// TODO: delete if not in use
func GrpcDial(endpoint string, opts []grpc.DialOption, runner Runner) {
	conn, err := grpc.Dial(endpoint, opts...)

	Check(err, "Could not initialise gRPC dial")
	defer conn.Close()

	runner(conn)
}
