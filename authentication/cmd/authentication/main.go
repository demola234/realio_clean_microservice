package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/demola234/authentication/config"
	db "github.com/demola234/authentication/db/sqlc"
	_ "github.com/demola234/authentication/docs/statik"
	pb "github.com/demola234/authentication/infrastructure/api/grpc"
	grpcHandler "github.com/demola234/authentication/infrastructure/api/user_handler"
	"github.com/demola234/authentication/internal/repository"
	usercase "github.com/demola234/authentication/internal/usecase"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	conn, err := sql.Open(configs.DBDriver, configs.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer conn.Close()

	dbQueries := db.New(conn)
	userRepo := repository.NewUserRepository(dbQueries)
	oAuthRepo := repository.NewOAuthRepository(&configs)
	userUsecase := usercase.NewUserUsecase(userRepo, oAuthRepo)

	server := grpcHandler.NewUserHandler(userUsecase)

	go runGRPCServer(configs, server)
	runGatewayServer(configs, server)
}

func runGRPCServer(configs config.Config, server pb.AuthServiceServer) {
	listener, err := net.Listen("tcp", configs.GRPCServerAddress)
	if err != nil {
		log.Fatalf("cannot start gRPC listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Printf("gRPC server running at %s", configs.GRPCServerAddress)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func runGatewayServer(configs config.Config, server pb.AuthServiceServer) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jsonOpt := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	mux := runtime.NewServeMux(jsonOpt)
	err := pb.RegisterAuthServiceHandlerServer(ctx, mux, server)
	if err != nil {
		log.Fatalf("cannot register gateway handler: %v", err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatalf("cannot create statik fs: %v", err)
	}
	httpMux.Handle("/swagger/", http.StripPrefix("/swagger", http.FileServer(statikFS)))

	listener, err := net.Listen("tcp", configs.HTTPServerAddress)
	if err != nil {
		log.Fatalf("cannot start HTTP listener: %v", err)
	}

	log.Printf("HTTP gateway server running at %s", configs.HTTPServerAddress)
	if err := http.Serve(listener, httpMux); err != nil {
		log.Fatalf("cannot serve HTTP gateway: %v", err)
	}
}
