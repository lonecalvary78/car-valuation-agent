package config

import "time"

type RateLimit struct {
	Limit  int
	Window time.Duration
}
