package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k1tasun/GoEdge-Gateway/api/proto"
	"github.com/k1tasun/GoEdge-Gateway/internal/config"
	"github.com/k1tasun/GoEdge-Gateway/internal/models"
	"github.com/k1tasun/GoEdge-Gateway/internal/mqtt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()
	slog.Info("starting mqtt ingestor", "broker", cfg.MQTTBroker)

	conn, err := grpc.NewClient("localhost:"+cfg.ServerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to connect to storage service", "error", err)
		os.Exit(1)
	}
	defer conn.Close()
	client := proto.NewStorageServiceClient(conn)

	ingestor, err := mqtt.NewIngestor(cfg.MQTTBroker, "goedge-ingestor", cfg.MQUITTopic)
	if err != nil {
		slog.Error("failed to create mqtt ingestor", "error", err)
		os.Exit(1)
	}
	defer ingestor.Close()

	if err := ingestor.Start(func(reading *models.SensorReading) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := client.StoreReading(ctx, &proto.StoreReadingRequest{
			Reading: &proto.SensorReading{
				DeviceId:  reading.DeviceID,
				Type:      reading.Type,
				Value:     reading.Value,
				Unit:      reading.Unit,
				Timestamp: timestamppb.New(reading.RecordedAt),
			},
		})
		if err != nil {
			slog.Error("error storing reading", "error", err)
		}
	}); err != nil {
		slog.Error("failed to start ingestor", "error", err)
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	slog.Info("shutting down mqtt ingestor")
}
