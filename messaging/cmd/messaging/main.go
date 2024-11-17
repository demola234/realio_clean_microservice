package main

import (
	"context"
	"log"
	"net"
	"time"

	"job_portal/messaging/config"
	"job_portal/messaging/db/mongo"
	pb "job_portal/messaging/interfaces/api/grpc"
	grpcHandler "job_portal/messaging/interfaces/api/message_handler"
	"job_portal/messaging/internal/repository"
	usercase "job_portal/messaging/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Initialize MongoDB
	client, err := mongo.NewClient(configs.MongoURI)
	if err != nil {
		log.Fatalf("failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("error disconnecting MongoDB client: %v", err)
		}
	}()

	db := client.Database(configs.DBUser)

	// Initialize repository and use case
	messageRepo := repository.NewMessageRepository(db, configs.DBUser)
	messageUsecase := usercase.NewMessagingUseCase(messageRepo)

	// Create the gRPC handler using the use case
	messageHandler := grpcHandler.NewMessageHandler(messageUsecase)

	// Start gRPC server
	lis, err := net.Listen("tcp", configs.GRPCServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessagingServiceServer(grpcServer, messageHandler)
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on %s", configs.GRPCServerAddress)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
