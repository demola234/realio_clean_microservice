package handler

import (
	"context"
	"io"
	"net/http"
	"strconv"

	errorResponse "github.com/demola234/api_gateway/infrastructure/error_response"
	"github.com/demola234/api_gateway/infrastructure/grpc_clients"
	token "github.com/demola234/api_gateway/infrastructure/middleware/token_maker"
	pb "github.com/demola234/authentication/infrastructure/api/grpc"

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

// OAuthLogin handles OAuth login requests
func (h *AuthHandler) OAuthLogin(c *gin.Context) {
	var req pb.OAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.AuthClient.Client.OAuthLogin(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) OAuthRegister(c *gin.Context) {
	var req pb.OAuthRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.AuthClient.Client.OAuthRegister(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) UploadImage(c *gin.Context) {
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}

	userID := authPayload.(*token.Payload).UserID

	image, err := c.FormFile("content")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read image file"})
		return
	}

	imageData, err := image.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open image file"})
		return
	}
	defer imageData.Close()

	fileBytes, err := io.ReadAll(imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image data"})
		return
	}

	res, err := h.AuthClient.Client.UploadImage(context.Background(), &pb.UploadImageRequest{
		UserId:  userID,
		Content: fileBytes,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ForgotPassword handles the forgot password request
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req pb.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.AuthClient.Client.ForgotPassword(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// VerifyResetPassword handles the verification of OTP for password reset
func (h *AuthHandler) VerifyResetPassword(c *gin.Context) {
	var req pb.VerifyResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.AuthClient.Client.VerifyResetPassword(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ResetPassword handles the password reset after OTP verification
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req pb.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	res, err := h.AuthClient.Client.ResetPassword(context.Background(), &req)
	if err != nil {
		if err.Error() == "You must verify the code before resetting your password" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ChangePassword handles changing the user's password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// Get user ID from authorization payload
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}
	userID := authPayload.(*token.Payload).UserID

	// Bind the request body
	var req pb.ChangePasswordRequest

	req.UserId = userID

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	// Create context with user ID for the gRPC service
	ctx := context.Background()

	// Call the gRPC service
	res, err := h.AuthClient.Client.ChangePassword(ctx, &req)
	if err != nil {
		// Check for specific error messages
		if err.Error() == "current password is incorrect" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
			return
		}
		if err.Error() == "new password does not meet requirements" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New password does not meet security requirements"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetProfile handles getting the user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from authorization payload
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}
	userID := authPayload.(*token.Payload).UserID

	// Create context
	ctx := context.Background()

	// Call the gRPC service
	res, err := h.AuthClient.Client.GetProfile(ctx, &pb.GetProfileRequest{
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateProfile handles updating the user's profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from authorization payload
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}
	userID := authPayload.(*token.Payload).UserID

	// Bind the request body
	var req pb.UpdateProfileRequest

	req.UserId = userID

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	// Create context
	ctx := context.Background()

	// Call the gRPC service
	res, err := h.AuthClient.Client.UpdateProfile(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetSessions handles getting all active sessions for the user
func (h *AuthHandler) GetSessions(c *gin.Context) {
	// Get user ID from authorization payload
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}
	userID := authPayload.(*token.Payload).UserID

	// Create context
	ctx := context.Background()

	// Call the gRPC service
	res, err := h.AuthClient.Client.GetSessions(ctx, &pb.GetSessionsRequest{
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// RevokeSession handles revoking a specific session
func (h *AuthHandler) RevokeSession(c *gin.Context) {
	// Get user ID from authorization payload
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}
	userID := authPayload.(*token.Payload).UserID

	// Get session ID from path
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session ID is required"})
		return
	}

	// Create context
	ctx := context.Background()

	// Call the gRPC service
	res, err := h.AuthClient.Client.RevokeSession(ctx, &pb.RevokeSessionRequest{
		SessionId: sessionID,
		UserId:    userID,
	})
	if err != nil {
		// Check for specific error messages
		if err.Error() == "unauthorized: session does not belong to the user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to revoke this session"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DeactivateAccount handles temporarily deactivating a user account
func (h *AuthHandler) DeactivateAccount(c *gin.Context) {
	// Get user ID from authorization payload
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}
	userID := authPayload.(*token.Payload).UserID

	// Bind the request body
	var req pb.DeactivateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	req.UserId = userID

	// Create context
	ctx := context.Background()

	// Call the gRPC service
	res, err := h.AuthClient.Client.DeactivateAccount(ctx, &req)
	if err != nil {
		// Check for specific error messages
		if err.Error() == "password is incorrect" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is incorrect"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DeleteAccount handles permanently deleting a user account
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	// Get user ID from authorization payload
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}
	userID := authPayload.(*token.Payload).UserID

	// Bind the request body
	var req pb.DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse.ErrInvalidRequest)
		return
	}

	req.UserId = userID

	// Create context
	ctx := context.Background()

	// Call the gRPC service
	res, err := h.AuthClient.Client.DeleteAccount(ctx, &req)
	if err != nil {
		// Check for specific error messages
		if err.Error() == "password is incorrect" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is incorrect"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetLoginHistory handles getting the login history for the user
func (h *AuthHandler) GetLoginHistory(c *gin.Context) {
	// Get user ID from authorization payload
	authPayload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization payload not found"})
		return
	}
	userID := authPayload.(*token.Payload).UserID

	// Parse limit query parameter with default value of 10
	limit := 10
	limitParam := c.DefaultQuery("limit", "10")
	if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
		limit = parsedLimit
	}

	// Create context
	ctx := context.Background()

	// Call the gRPC service
	res, err := h.AuthClient.Client.GetLoginHistory(ctx, &pb.GetLoginHistoryRequest{
		Limit:  int32(limit),
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
