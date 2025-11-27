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
	query := `INSERT INTO readings (device_id, reading_type, value, unit, recorded_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, reading.DeviceID, reading.Type, reading.Value, reading.Unit, reading.RecordedAt).Scan(&reading.ID, &reading.CreatedAt)
}

func (r *Repository) SaveBatch(ctx context.Context, readings []*models.SensorReading) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO readings (device_id, reading_type, value, unit, recorded_at) VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	count := 0
	for _, reading := range readings {
		if _, err := stmt.ExecContext(ctx, reading.DeviceID, reading.Type, reading.Value, reading.Unit, reading.RecordedAt); err != nil {
			return count, err
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) GetReadingsByDevice(ctx context.Context, deviceID string, limit int) ([]*models.SensorReading, error) {
	query := `SELECT id, device_id, reading_type, value, unit, recorded_at, created_at FROM readings WHERE device_id = $1 ORDER BY recorded_at DESC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, deviceID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch readings: %w", err)
	}
	defer rows.Close()

	var readings []*models.SensorReading
	for rows.Next() {
		var r models.SensorReading
		if err := rows.Scan(&r.ID, &r.DeviceID, &r.Type, &r.Value, &r.Unit, &r.RecordedAt, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan reading: %w", err)
		}
		readings = append(readings, &r)
	}
	return readings, rows.Err()
}
