package routes

import (
	"github.com/demola234/api_gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterMessageRoutes(rg *gin.RouterGroup, messageHandler *handler.MessageHandler, authMiddleware gin.HandlerFunc) {
	messageRoutes := rg.Group("/message")

	{
		messageRoutes.GET("/", authMiddleware, messageHandler.GetMessages)
		messageRoutes.GET("/conversation", authMiddleware, messageHandler.GetConversationBetweenUser)
		messageRoutes.GET("/:id", authMiddleware, messageHandler.GetConversationByID)
		messageRoutes.POST("/", authMiddleware, messageHandler.SendMessage)
	}
}
