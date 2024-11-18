package handler

import (
	"context"
	"net/http"

	errorResponse "job_portal/api_gateway/infrastructure/error_response"
	"job_portal/api_gateway/infrastructure/grpc_clients"
	pb "job_portal/messaging/infrastructure/api/grpc"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	MessageClient *grpc_clients.MessageClient
}

func NewMessageHandler(messageClient *grpc_clients.MessageClient) *MessageHandler {
	return &MessageHandler{MessageClient: messageClient}
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	var req pb.GetMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.MessageClient.Client.GetMessages(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req pb.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.MessageClient.Client.SendMessage(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

func (h *MessageHandler) GetConversationBetweenUser(c *gin.Context) {
	var req pb.GetConversationBetweenUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.MessageClient.Client.GetConversationBetweenUsers(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *MessageHandler) GetConversationByID(c *gin.Context) {
	var req pb.GetConversationBetweenUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.MessageClient.Client.GetConversationBetweenUsers(context.Background(), &pb.GetConversationBetweenUsersRequest{User1Id: req.User1Id, User2Id: req.User2Id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}
