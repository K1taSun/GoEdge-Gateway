package models

import "time"

type SensorReading struct {
	ID          int64     `json:"id"`
	DeviceID    string    `json:"device_id"`
	Type        string    `json:"type"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	RecordedAt  time.Time `json:"recorded_at"`
	CreatedAt   time.Time `json:"created_at"`
}

