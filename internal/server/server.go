package server

import (
	"context"
	"log/slog"

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
	err := s.repo.SaveReading(ctx, &models.SensorReading{
		DeviceID:   req.Reading.DeviceId,
		Type:       req.Reading.Type,
		Value:      req.Reading.Value,
		Unit:       req.Reading.Unit,
		RecordedAt: req.Reading.Timestamp.AsTime(),
	})
	if err != nil {
		slog.Error("failed to save reading", "error", err)
		return &pb.StoreReadingResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.StoreReadingResponse{Success: true, Message: "Stored successfully"}, nil
}

func (s *GatewayServer) StoreBatch(ctx context.Context, req *pb.StoreBatchRequest) (*pb.StoreBatchResponse, error) {
	readings := make([]*models.SensorReading, len(req.Readings))
	for i, r := range req.Readings {
		readings[i] = &models.SensorReading{
			DeviceID:   r.DeviceId,
			Type:       r.Type,
			Value:      r.Value,
			Unit:       r.Unit,
			RecordedAt: r.Timestamp.AsTime(),
		}
	}

	count, err := s.repo.SaveBatch(ctx, readings)
	if err != nil {
		slog.Error("failed to save batch", "error", err)
		return &pb.StoreBatchResponse{Success: false}, nil
	}
	return &pb.StoreBatchResponse{Success: true, Count: int32(count)}, nil
}

func (s *GatewayServer) GetReadings(ctx context.Context, req *pb.GetReadingsRequest) (*pb.GetReadingsResponse, error) {
	readings, err := s.repo.GetReadingsByDevice(ctx, req.DeviceId, int(req.Limit))
	if err != nil {
		slog.Error("failed to get readings", "error", err)
		return nil, err
	}

	pbReadings := make([]*pb.SensorReading, len(readings))
	for i, r := range readings {
		pbReadings[i] = &pb.SensorReading{
			DeviceId:  r.DeviceID,
			Type:      r.Type,
			Value:     r.Value,
			Unit:      r.Unit,
			Timestamp: timestamppb.New(r.RecordedAt),
		}
	}

	return &pb.GetReadingsResponse{Readings: pbReadings}, nil
}
