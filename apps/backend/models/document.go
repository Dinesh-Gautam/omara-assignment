package models

import "time"

type Document struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	FileName        string    `json:"file_name"`
	GCSPath         string    `json:"gcs_path"`
	Status          string    `json:"status"` // e.g., "processing", "processed", "failed"
	ProcessingError string    `json:"processingError,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
