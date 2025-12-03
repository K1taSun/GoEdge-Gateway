package main

import (
	"context"
	"log"
	"time"

	"github.com/k1tasun/GoEdge-Gateway/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// Connect to the storage service
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewStorageServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deviceID := "e2e-test-device"

	// 1. Store a reading
	log.Println("Storing reading...")
	_, err = client.StoreReading(ctx, &proto.StoreReadingRequest{
		Reading: &proto.SensorReading{
			DeviceId:  deviceID,
			Type:      "temperature",
			Value:     23.5,
			Unit:      "Celsius",
			Timestamp: timestamppb.Now(),
		},
	})
	if err != nil {
		log.Fatalf("could not store reading: %v", err)
	}
	log.Println("Reading stored successfully")

	// Allow some time for processing if needed (though gRPC is synchronous here)
	time.Sleep(1 * time.Second)

	// 2. Get readings
	log.Println("Fetching readings...")
	resp, err := client.GetReadings(ctx, &proto.GetReadingsRequest{
		DeviceId: deviceID,
		Limit:    10,
	})
	if err != nil {
		log.Fatalf("could not get readings: %v", err)
	}

	// 3. Verify
	if len(resp.Readings) == 0 {
		log.Fatal("Expected at least one reading, got 0")
	}

	found := false
	for _, r := range resp.Readings {
		if r.DeviceId == deviceID && r.Value == 23.5 {
			found = true
			log.Printf("Found reading: DeviceID=%s, Value=%f, Type=%s", r.DeviceId, r.Value, r.Type)
			break
		}
	}

	if found {
		log.Println("E2E Test PASSED: Successfully stored and retrieved reading.")
	} else {
		log.Fatal("E2E Test FAILED: Stored reading not found in response.")
	}
}
