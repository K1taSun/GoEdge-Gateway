package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/k1tasun/GoEdge-Gateway/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveReading(ctx context.Context, reading *models.SensorReading) error {
	query := `
		INSERT INTO readings (device_id, reading_type, value, unit, recorded_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query,
		reading.DeviceID,
		reading.Type,
		reading.Value,
		reading.Unit,
		reading.RecordedAt,
	).Scan(&reading.ID, &reading.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to save reading: %w", err)
	}

	return nil
}

func (r *Repository) GetReadingsByDevice(ctx context.Context, deviceID string, limit int) ([]*models.SensorReading, error) {
	query := `
		SELECT id, device_id, reading_type, value, unit, recorded_at, created_at
		FROM readings
		WHERE device_id = $1
		ORDER BY recorded_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, deviceID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch readings: %w", err)
	}
	defer rows.Close()

	var readings []*models.SensorReading
	for rows.Next() {
		var r models.SensorReading
		if err := rows.Scan(
			&r.ID, &r.DeviceID, &r.Type, &r.Value, &r.Unit, &r.RecordedAt, &r.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan reading: %w", err)
		}
		readings = append(readings, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return readings, nil
}
