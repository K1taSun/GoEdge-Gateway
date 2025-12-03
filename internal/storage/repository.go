package storage

import (
	"context"

	"github.com/k1tasun/GoEdge-Gateway/internal/models"
)

type Repository interface {
	SaveReading(ctx context.Context, reading *models.SensorReading) error
	SaveBatch(ctx context.Context, readings []*models.SensorReading) (int, error)
	GetReadingsByDevice(ctx context.Context, deviceID string, limit int) ([]*models.SensorReading, error)
}
