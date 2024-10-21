package routes

import (
	"job_portal/api_gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler) {
	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}
}
