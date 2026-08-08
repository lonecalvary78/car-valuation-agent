package model

import (
	"time"
)

type HealthCheckResponse struct {
	Status    string    `json:"status"`
	CheckedAt time.Time `json:"checked_at"`
}

func OfResponse(status string) HealthCheckResponse {
	return HealthCheckResponse{
		Status:    status,
		CheckedAt: time.Now(),
	}
}
