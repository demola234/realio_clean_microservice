package main

import (
	"database/sql"
	"log"
	"net"

	"job_portal/authentication/config"
	db "job_portal/authentication/db/sqlc"
	grpcHandler "job_portal/authentication/interfaces/api/user_handler"
	pb "job_portal/authentication/interfaces/api/grpc"
	"job_portal/authentication/internal/repository"
	usercase "job_portal/authentication/internal/usecase"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

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

	// Start gRPC server
	lis, err := net.Listen("tcp", configs.GRPCServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on %s", configs.GRPCServerAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
