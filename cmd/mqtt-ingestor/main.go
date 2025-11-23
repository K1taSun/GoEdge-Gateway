package main

import (
	"context"
	"log"
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
	cfg := config.Load()

	log.Printf("Starting MQTT Ingestor connected to %s", cfg.MQTTBroker)

	// Connect to gRPC Storage Service
	conn, err := grpc.NewClient("localhost:"+cfg.ServerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewStorageServiceClient(conn)

	// Initialize MQTT Ingestor
	ingestor, err := mqtt.NewIngestor(cfg.MQTTBroker, "goedge-ingestor", cfg.MQUITTopic)
	if err != nil {
		log.Fatalf("Failed to create MQTT ingestor: %v", err)
	}
	defer ingestor.Close()

	// Handle incoming readings
	handler := func(reading *models.SensorReading) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		pbReading := &proto.SensorReading{
			DeviceId:  reading.DeviceID,
			Type:      reading.Type,
			Value:     reading.Value,
			Unit:      reading.Unit,
			Timestamp: timestamppb.New(reading.RecordedAt),
		}

		_, err := client.StoreReading(ctx, &proto.StoreReadingRequest{Reading: pbReading})
		if err != nil {
			log.Printf("Error storing reading: %v", err)
		}
	}

	if err := ingestor.Start(handler); err != nil {
		log.Fatalf("Failed to start ingestor: %v", err)
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down MQTT Ingestor...")
}
