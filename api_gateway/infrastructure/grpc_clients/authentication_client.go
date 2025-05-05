package grpc_clients

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/demola234/authentication/infrastructure/api/grpc"

	"google.golang.org/grpc"
)

type AuthenticationClient struct {
	Client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

// NewAuthenticationClient creates a new gRPC client for the Authentication service
func NewAuthenticationClient(address string, timeout time.Duration) (*AuthenticationClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Establish the connection with timeout
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", address, err)
	}

	log.Printf("Connected to Authentication service at %s", address)
	client := pb.NewAuthServiceClient(conn)

	return &AuthenticationClient{
		Client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (ac *AuthenticationClient) Close() error {
	if err := ac.conn.Close(); err != nil {
		return fmt.Errorf("failed to close gRPC connection: %w", err)
	}
	log.Println("gRPC connection closed")
	return nil
}
