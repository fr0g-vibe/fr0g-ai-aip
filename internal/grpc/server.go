package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// StartGRPCServer starts the gRPC server on the specified port
// Note: Full gRPC implementation requires protobuf generation
func StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	
	// Enable reflection for debugging with tools like grpcurl
	reflection.Register(s)

	fmt.Printf("gRPC server listening on port %s (basic server - run 'make install-proto-tools && make proto' for full functionality)\n", port)
	return s.Serve(lis)
}
