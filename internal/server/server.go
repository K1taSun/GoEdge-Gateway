package server

import (
	"context"
	"log"

	pb "github.com/k1tasun/GoEdge-Gateway/api/proto"
	"github.com/k1tasun/GoEdge-Gateway/internal/models"
	"github.com/k1tasun/GoEdge-Gateway/internal/storage/postgres"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GatewayServer struct {
	pb.UnimplementedStorageServiceServer
	repo *postgres.Repository
}

func NewGatewayServer(repo *postgres.Repository) *GatewayServer {
	return &GatewayServer{repo: repo}
}

func (s *GatewayServer) StoreReading(ctx context.Context, req *pb.StoreReadingRequest) (*pb.StoreReadingResponse, error) {
	reading := req.Reading
	model := &models.SensorReading{
		DeviceID:   reading.DeviceId,
		Type:       reading.Type,
		Value:      reading.Value,
		Unit:       reading.Unit,
		RecordedAt: reading.Timestamp.AsTime(),
	}

	if err := s.repo.SaveReading(ctx, model); err != nil {
		log.Printf("Failed to save reading: %v", err)
		return &pb.StoreReadingResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.StoreReadingResponse{Success: true, Message: "Stored successfully"}, nil
}

func (s *GatewayServer) StoreBatch(ctx context.Context, req *pb.StoreBatchRequest) (*pb.StoreBatchResponse, error) {
	count := 0
	for _, r := range req.Readings {
		model := &models.SensorReading{
			DeviceID:   r.DeviceId,
			Type:       r.Type,
			Value:      r.Value,
			Unit:       r.Unit,
			RecordedAt: r.Timestamp.AsTime(),
		}
		if err := s.repo.SaveReading(ctx, model); err == nil {
			count++
		}
	}
	return &pb.StoreBatchResponse{Success: true, Count: int32(count)}, nil
}

func ConvertToProto(m *models.SensorReading) *pb.SensorReading {
	return &pb.SensorReading{
		DeviceId:  m.DeviceID,
		Type:      m.Type,
		Value:     m.Value,
		Unit:      m.Unit,
		Timestamp: timestamppb.New(m.RecordedAt),
	}
}

