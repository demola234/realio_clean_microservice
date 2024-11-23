package routes

import (
	"job_portal/api_gateway/internal/handler"
	// Import middleware package

	"github.com/gin-gonic/gin"
)

func RegisterPropertyRoutes(rg *gin.RouterGroup, propertyHandler *handler.PropertyHandler, authMiddleware gin.HandlerFunc) {
	propertyRoutes := rg.Group("/property")

	{

		propertyRoutes.GET("/", authMiddleware, propertyHandler.GetProperties)
		propertyRoutes.GET("/user", authMiddleware, propertyHandler.GetPropertiesByOwner)
		propertyRoutes.GET("/:id", authMiddleware, propertyHandler.GetProperty)    // GET /properties/:id
		propertyRoutes.POST("/", authMiddleware, propertyHandler.CreateProperty)   // POST /properties
		propertyRoutes.PUT("/:id", authMiddleware, propertyHandler.UpdateProperty) // PUT /properties/:id
	}
}
