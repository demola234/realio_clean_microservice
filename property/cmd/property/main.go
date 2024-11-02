package main

import (
	"database/sql"
	"job_portal/property/config"
	db "job_portal/property/db/sqlc"
	pb "job_portal/property/interfaces/api/grpc"
	grpcHandler "job_portal/property/interfaces/api/property_handler"
	"job_portal/property/internal/repository"
	"job_portal/property/internal/usecases"
	"log"
	"net"

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
	propertyRepo := repository.NewPropertyRepository(dbQueries)
	propertyUsecase := usecases.NewPropertyUsecase(propertyRepo)

	propertyService := grpcHandler.NewPropertyHandler(propertyUsecase)

	// Start the gRPC server
	lis, err := net.Listen("tcp", configs.GRPCServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on %s", configs.GRPCServerAddress)

	grpcServer := grpc.NewServer()
	pb.RegisterPropertyServiceServer(grpcServer, propertyService)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
