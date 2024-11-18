package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"job_portal/messaging/config"
	"job_portal/messaging/db/mongo"
	"job_portal/messaging/infrastructure/socket"
	"job_portal/messaging/internal/repository"
	"job_portal/messaging/internal/usecase"
)

func main() {
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

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
	messageRepo := repository.NewMessageRepository(db, "messages")
	conversationRepo := repository.NewConversationRepository(db, "conversations")
	messageUsecase := usecase.NewMessagingUseCase(messageRepo, conversationRepo)

	webSocketHandler := socket.NewMessageWebSocket(messageUsecase)

	http.HandleFunc("/ws", webSocketHandler.HandleWebSocket)
	log.Println("WebSocket server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
