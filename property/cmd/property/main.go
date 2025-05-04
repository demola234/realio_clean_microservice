package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/demola234/property/config"
	db "github.com/demola234/property/db/sqlc"
	pb "github.com/demola234/property/infrastructure/api/grpc"
	grpcHandler "github.com/demola234/property/infrastructure/api/property_handler"
	"github.com/demola234/property/infrastructure/messaging/kafka"
	"github.com/demola234/property/internal/repository"
	"github.com/demola234/property/internal/usecases"

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

	// Initialize Kafka producer
	kafkaProducer := kafka.NewKafkaProducer(configs.KafkaBrokers, configs.KafkaGroupID)
	defer kafkaProducer.Close()

	// Pass the Kafka producer to the use case
	propertyUsecase := usecases.NewPropertyUsecase(propertyRepo, kafkaProducer)

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
