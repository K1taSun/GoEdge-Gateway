package main

import (
	"log/slog"
	"net"
	"os"

	"github.com/k1tasun/GoEdge-Gateway/api/proto"
	"github.com/k1tasun/GoEdge-Gateway/internal/config"
	"github.com/k1tasun/GoEdge-Gateway/internal/server"
	"github.com/k1tasun/GoEdge-Gateway/internal/storage/postgres"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()
	slog.Info("starting storage service", "port", cfg.ServerPort)

	db, err := postgres.NewConnection(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterStorageServiceServer(grpcServer, server.NewGatewayServer(postgres.NewRepository(db)))

	slog.Info("grpc server listening", "address", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}
