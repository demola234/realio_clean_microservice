package handler

import (
	"context"
	"net/http"

	errorResponse "job_portal/api_gateway/interfaces/error_response"
	"job_portal/api_gateway/interfaces/grpc_clients"
	pb "job_portal/authentication/interfaces/api/grpc"
	token "job_portal/api_gateway/interfaces/middleware/token_maker"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.AuthClient.Client.Login(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse.ErrInvalidRequest)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) GetUser(c *gin.Context) {
	// Get userID from authorization payload in context
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	userID := authPayload.(*token.Payload).UserID

	res, err := h.AuthClient.Client.GetUser(context.Background(), &pb.GetUserRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) VerifyUser(c *gin.Context) {
	var req pb.VerifyUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.AuthClient.Client.VerifyUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) ResendOtp(c *gin.Context) {
	var req pb.ResendOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.AuthClient.Client.ResendOtp(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Logout(c *gin.Context) {

	// Get userID from authorization payload in context
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	userID := authPayload.(*token.Payload)

	req := pb.LogOutRequest{
		UserId: userID.UserID,
	}

	res, err := h.AuthClient.Client.LogOut(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
