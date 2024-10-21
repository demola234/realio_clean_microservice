package main

import (
	"job_portal/api_gateway/config"
	"job_portal/api_gateway/interfaces/grpc_clients"
	"job_portal/api_gateway/internal/handler"
	auth "job_portal/api_gateway/routes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}

	// Initialize gRPC client
	authClient, err := grpc_clients.NewAuthenticationClient(configs.AuthenticationGRPCAddress)
	if err != nil {
		log.Fatalf("Failed to connect to Authentication service: %v", err)
	}

	// Create a new Gin router
	router := gin.Default()

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authClient)

	// Group routes under /v1
	v1 := router.Group("/v1")
	{
		// Health check endpoint
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "OK",
			})
		})
	}

	// Define authentication routes
	auth.RegisterRoutes(v1, authHandler)

	// Start the HTTP server
	log.Printf("Starting API Gateway at %s...", configs.Port)
	if err := router.Run(configs.Port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
