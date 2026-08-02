package model

import (
	"time"

	"github.com/bytedance/sonic"
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

func (hc *HealthCheckResponse) ToJSON() ([]byte, error) {
	return sonic.Marshal(hc)
}
