package main

import (
	"log"
	"net"

	"github.com/k1tasun/GoEdge-Gateway/api/proto"
	"github.com/k1tasun/GoEdge-Gateway/internal/config"
	"github.com/k1tasun/GoEdge-Gateway/internal/server"
	"github.com/k1tasun/GoEdge-Gateway/internal/storage/postgres"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting Storage Service on port %s", cfg.ServerPort)

	// Initialize Database
	db, err := postgres.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewRepository(db)

	// Setup gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	gatewayServer := server.NewGatewayServer(repo)
	proto.RegisterStorageServiceServer(grpcServer, gatewayServer)

	log.Printf("gRPC server listening on :%s", cfg.ServerPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
