package socket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/demola234/messaging/internal/domain/entity"
	"github.com/demola234/messaging/internal/usecase"

	"github.com/gorilla/websocket"
)

type MessageWebSocket struct {
	upgrader    websocket.Upgrader
	messagingUC usecase.MessagingUseCase
	clients     map[*websocket.Conn]string
}

func NewMessageWebSocket(messagingUC usecase.MessagingUseCase) *MessageWebSocket {
	return &MessageWebSocket{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		messagingUC: messagingUC,
		clients:     make(map[*websocket.Conn]string),
	}
}

func (ws *MessageWebSocket) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("New WebSocket connection: %s", conn.RemoteAddr().String())

	for {
		var data map[string]string
		err := conn.ReadJSON(&data)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(ws.clients, conn)
			break
		}

		event := data["event"]
		switch event {
		case "join":
			ws.handleJoin(conn, data)
		case "send_message":
			ws.handleSendMessage(conn, data)
		default:
			conn.WriteJSON(map[string]string{"error": "Unknown event"})
		}
	}
}

func (ws *MessageWebSocket) handleJoin(conn *websocket.Conn, data map[string]string) {
	userID, ok := data["userId"]
	if !ok || userID == "" {
		conn.WriteJSON(map[string]string{"error": "Missing or invalid userId"})
		return
	}

	ws.clients[conn] = userID
	log.Printf("User %s joined", userID)
	conn.WriteJSON(map[string]string{"status": "success", "room": userID})
}

func (ws *MessageWebSocket) handleSendMessage(conn *websocket.Conn, data map[string]string) {
	if err := validateMessagePayload(data); err != nil {
		log.Printf("Invalid message payload: %+v, Error: %v", data, err)
		conn.WriteJSON(map[string]string{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := &entity.Message{
		ConversationID: data["conversationId"],
		SenderID:       data["senderId"],
		ReceiverID:     data["receiverId"],
		Content:        data["content"],
	}

	savedMessage, err := ws.messagingUC.SendMessage(ctx, message)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		conn.WriteJSON(map[string]string{"error": "Failed to send message"})
		return
	}

	conn.WriteJSON(savedMessage)

	for client, userID := range ws.clients {
		if userID == savedMessage.ReceiverID {
			client.WriteJSON(map[string]interface{}{
				"event":   "new_message",
				"message": savedMessage,
			})
		}
	}

	log.Printf("Message sent: %+v", savedMessage)
}

func validateMessagePayload(data map[string]string) error {
	requiredFields := []string{"senderId", "receiverId", "content"}
	for _, field := range requiredFields {
		if value, exists := data[field]; !exists || value == "" {
			return fmt.Errorf("missing or invalid field: %s", field)
		}
	}
	return nil
}
