package routes

import (
	"job_portal/api_gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterMessageRoutes(rg *gin.RouterGroup, messageHandler *handler.MessageHandler, authMiddleware gin.HandlerFunc) {
	messageRoutes := rg.Group("/message")

	{
		messageRoutes.GET("/", authMiddleware, messageHandler.GetMessages)                            // GET /message
		messageRoutes.GET("/conversation", authMiddleware, messageHandler.GetConversationBetweenUser) // GET /message/conversation
		messageRoutes.GET("/:id", authMiddleware, messageHandler.GetConversationByID)                 // GET /message/:id
		messageRoutes.POST("/", authMiddleware, messageHandler.SendMessage)                           // POST /message
	}
}
