package handler

import (
	"context"
	"net/http"

	"job_portal/api_gateway/interfaces/grpc_clients"
	errorResponse "job_portal/api_gateway/interfaces/error_response"
	pb "job_portal/authentication/interfaces/api/grpc"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthClient *grpc_clients.AuthenticationClient
}



func NewAuthHandler(authClient *grpc_clients.AuthenticationClient) *AuthHandler {
	return &AuthHandler{AuthClient: authClient}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req pb.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.AuthClient.Client.Register(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req pb.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.AuthClient.Client.Login(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse.ErrInvalidRequest)
		return
	}

	c.JSON(http.StatusOK, res)
}
