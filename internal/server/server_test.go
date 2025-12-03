package server

import (
	"context"
	"testing"
	"time"

	"github.com/k1tasun/GoEdge-Gateway/api/proto"
	"github.com/k1tasun/GoEdge-Gateway/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MockRepository is a mock implementation of storage.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveReading(ctx context.Context, reading *models.SensorReading) error {
	args := m.Called(ctx, reading)
	return args.Error(0)
}

func (m *MockRepository) SaveBatch(ctx context.Context, readings []*models.SensorReading) (int, error) {
	args := m.Called(ctx, readings)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetReadingsByDevice(ctx context.Context, deviceID string, limit int) ([]*models.SensorReading, error) {
	args := m.Called(ctx, deviceID, limit)
	return args.Get(0).([]*models.SensorReading), args.Error(1)
}

func TestStoreReading(t *testing.T) {
	mockRepo := new(MockRepository)
	server := NewGatewayServer(mockRepo)
	ctx := context.Background()

	req := &proto.StoreReadingRequest{
		Reading: &proto.SensorReading{
			DeviceId:  "device-1",
			Type:      "temp",
			Value:     25.5,
			Unit:      "C",
			Timestamp: timestamppb.Now(),
		},
	}

	mockRepo.On("SaveReading", ctx, mock.AnythingOfType("*models.SensorReading")).Return(nil)

	resp, err := server.StoreReading(ctx, req)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	mockRepo.AssertExpectations(t)
}

func TestGetReadings(t *testing.T) {
	mockRepo := new(MockRepository)
	server := NewGatewayServer(mockRepo)
	ctx := context.Background()

	now := time.Now()
	readings := []*models.SensorReading{
		{
			DeviceID:   "device-1",
			Type:       "temp",
			Value:      25.5,
			Unit:       "C",
			RecordedAt: now,
		},
	}

	mockRepo.On("GetReadingsByDevice", ctx, "device-1", 10).Return(readings, nil)

	req := &proto.GetReadingsRequest{
		DeviceId: "device-1",
		Limit:    10,
	}

	resp, err := server.GetReadings(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, resp.Readings, 1)
	assert.Equal(t, 25.5, resp.Readings[0].Value)
	mockRepo.AssertExpectations(t)
}
