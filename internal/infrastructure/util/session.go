package util

import "github.com/google/uuid"

func NewSessionId() uuid.UUID {
	return uuid.New()
}
