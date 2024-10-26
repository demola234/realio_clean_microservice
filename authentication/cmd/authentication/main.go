package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	zapp "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"job_portal/authentication/config"
	db "job_portal/authentication/db/sqlc"
	pb "job_portal/authentication/interfaces/api/grpc"
	grpcHandler "job_portal/authentication/interfaces/api/user_handler"
	"job_portal/authentication/internal/repository"
	usercase "job_portal/authentication/internal/usecase"

	_ "github.com/lib/pq"
)

// Initialize logger using zap
func initLogger() *zap.Logger {
	logger, err := zap.NewProduction() // Use NewDevelopment() for development mode
	if err != nil {
		log.Fatalf("Failed to create zap logger: %v", err)
	}
	return logger
}

// Define a basic recovery handler
func grpcRecoveryHandler(p interface{}) (err error) {
	log.Printf("Recovered from panic: %v", p)
	return status.Errorf(codes.Internal, "Recovered from panic: %v", p)
}

// StartMetricsServer starts a simple HTTP server to expose Prometheus metrics
func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler()) // Expose the /metrics endpoint
	log.Println("Prometheus metrics available at http://localhost:9091/metrics")

	// Start HTTP server on port 8080
	if err := http.ListenAndServe(":9091", nil); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}

// Main function
func main() {
	// Load configuration
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Initialize logger
	logger := initLogger()
	defer logger.Sync() // Flushes buffer, if any

	// Connect to the database
	conn, err := sql.Open(configs.DBDriver, configs.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing database connection: %v", err)
		}
	}()

	// Initialize repository and use case
	dbQueries := db.New(conn)
	userRepo := repository.NewUserRepository(dbQueries)
	userUsecase := usercase.NewUserUsecase(userRepo)

	// Create the gRPC handler using the use case
	authService := grpcHandler.NewUserHandler(userUsecase)

	// Start gRPC server listener
	lis, err := net.Listen("tcp", configs.GRPCServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create the gRPC server with interceptors
	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_prometheus.UnaryServerInterceptor,
			zapp.UnaryServerInterceptor(logger),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcRecoveryHandler)),
		),
	)

	// Register the authentication service
	pb.RegisterAuthServiceServer(grpcServer, authService)
	grpc_prometheus.Register(grpcServer) // Register Prometheus metrics with the gRPC server
	reflection.Register(grpcServer)

	// Start the Prometheus metrics server in a separate goroutine
	go startMetricsServer()

	log.Printf("gRPC server is running on %s", configs.GRPCServerAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
