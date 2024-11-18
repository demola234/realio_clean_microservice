package grpc_clients

import (
	"context"
	"fmt"
	pb "job_portal/messaging/infrastructure/api/grpc"
	"log"
	"time"

	"google.golang.org/grpc"
)

type MessageClient struct {
	Client pb.MessagingServiceClient
	conn   *grpc.ClientConn
}

// NewAuthenticationClient creates a new gRPC client for the Authentication service
func NewMessagingClient(address string, timeout time.Duration) (*MessageClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Establish the connection with timeout
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", address, err)
	}

	log.Printf("Connected to Messaging service at %s", address)
	client := pb.NewMessagingServiceClient(conn)

	return &MessageClient{
		Client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (ac *MessageClient) Close() error {
	if err := ac.conn.Close(); err != nil {
		return fmt.Errorf("failed to close gRPC connection: %w", err)
	}
	log.Println("gRPC connection closed")
	return nil
}
