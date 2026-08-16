package model

import (
	"car-valuation-agent/internal/infrastructure/validator"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/bytedance/sonic"
)

type ChatRequest struct {
	SessionId string `json:"sessionId"`
	UserId    string `json:"userId"`
	Message   string `json:"message"`
}

func (chat *ChatRequest) FromJSON(requestBody []byte) error {
	err := sonic.Unmarshal(requestBody, chat)
	if err != nil {
		return fmt.Errorf("model: failed to unmarshal chat request: %w", err)
	}
	return nil
}

func (chat *ChatRequest) Validate() error {
	validationErrors := make([]error, 0)
	validationError := validator.ValidateForRequiredOfString("UserId", chat.UserId)
	if validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	return errors.Join(validationErrors...)
}

func (chat *ChatRequest) SetSessionId(sesionId uuid.UUID) {
	chat.SessionId = sesionId.String()
}

func (chat *ChatRequest) SetUserId(userId uuid.UUID) {
	chat.UserId = userId.String()
}
