package grpc_clients

import (
	"fmt"
	pb "job_portal/authentication/interfaces/api/grpc"
	"log"

	"google.golang.org/grpc"
)

type AuthenticationClient struct {
	Client pb.AuthServiceClient
}

func NewAuthenticationClient(address string) (*AuthenticationClient, error) {
	// Create a context with a timeout to avoid indefinite dialing attempts.

	// Use grpc.DialContext with context and a proper timeout.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	log.Printf("Connected to Authentication service at %s", address)
	client := pb.NewAuthServiceClient(conn)

	return &AuthenticationClient{Client: client}, nil
}
