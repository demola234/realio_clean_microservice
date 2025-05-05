package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/demola234/messaging/config"
	"github.com/demola234/messaging/db/mongo"
	pb "github.com/demola234/messaging/infrastructure/api/grpc"
	grpcHandler "github.com/demola234/messaging/infrastructure/api/message_handler"
	"github.com/demola234/messaging/infrastructure/socket"
	"github.com/demola234/messaging/internal/repository"
	"github.com/demola234/messaging/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	// Initialize MongoDB client
	client, err := mongo.NewClient(configs.MongoURI)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB client: %v", err)
		}
	}()

	// Repositories and UseCase initialization
	db := client.Database(configs.DBUser)
	messageRepo := repository.NewMessageRepository(db, "messages")
	conversationRepo := repository.NewConversationRepository(db, "conversations")
	messageUsecase := usecase.NewMessagingUseCase(messageRepo, conversationRepo)

	// WebSocket handler
	webSocketHandler := socket.NewMessageWebSocket(messageUsecase)

	// Start gRPC and HTTP servers in separate goroutines
	errChan := make(chan error)

	go func() {
		if err := startGRPCServer(configs.GRPCServerAddress, messageUsecase); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := startWebSocketServer(":8080", webSocketHandler); err != nil {
			errChan <- err
		}
	}()

	// Gracefully handle shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	select {
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	case <-signalChan:
		log.Println("Shutting down gracefully...")
		cancel()
	}
}

func startGRPCServer(address string, messageUsecase usecase.MessagingUseCase) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	messageService := grpcHandler.NewMessageHandler(messageUsecase)

	pb.RegisterMessagingServiceServer(grpcServer, messageService)
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on %s", address)
	return grpcServer.Serve(lis)
}

func startWebSocketServer(address string, handler *socket.MessageWebSocket) error {
	http.HandleFunc("/ws", handler.HandleWebSocket)

	server := &http.Server{
		Addr:         address,
		Handler:      nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("WebSocket server is running on %s", address)
	return server.ListenAndServe()
}
