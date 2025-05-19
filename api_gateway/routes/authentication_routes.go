package routes

import (
	"github.com/demola234/api_gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler, authMiddleware gin.HandlerFunc) {
	authRoutes := rg.Group("/auth")

	{
		// Registration and login
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/verify", authHandler.VerifyUser)
		authRoutes.POST("/resend-otp", authHandler.ResendOtp)

		// OAuth routes
		authRoutes.POST("/register-oauth", authHandler.OAuthRegister)
		authRoutes.POST("/login-oauth", authHandler.OAuthLogin)

		// Password reset flow
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/verify-reset", authHandler.VerifyResetPassword)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
	}

	// Protected routes (require authentication)
	{
		// User information
		authRoutes.GET("/user", authMiddleware, authHandler.GetUser)
		authRoutes.GET("/profile", authMiddleware, authHandler.GetProfile)
		authRoutes.PUT("/profile", authMiddleware, authHandler.UpdateProfile)

		// Session management
		authRoutes.POST("/logout", authMiddleware, authHandler.Logout)
		authRoutes.GET("/sessions", authMiddleware, authHandler.GetSessions)
		authRoutes.DELETE("/sessions/:session_id", authMiddleware, authHandler.RevokeSession)

		// Account management
		authRoutes.POST("/change-password", authMiddleware, authHandler.ChangePassword)
		authRoutes.POST("/upload-image", authMiddleware, authHandler.UploadImage)
		authRoutes.POST("/account/deactivate", authMiddleware, authHandler.DeactivateAccount)
		authRoutes.DELETE("/account", authMiddleware, authHandler.DeleteAccount)
		authRoutes.GET("/account/login-history", authMiddleware, authHandler.GetLoginHistory)
	}
}
