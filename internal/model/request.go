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
	// UserId is never taken from client input; it is set by the handler from the
	// authenticated request context so a caller cannot act on another user's behalf.
	UserId  string `json:"-"`
	Message string `json:"message"`
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
	validationError := validator.ValidateForRequiredOfString("Message", chat.Message)
	if validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	return errors.Join(validationErrors...)
}

func (chat *ChatRequest) SetSessionId(sesionId uuid.UUID) {
	chat.SessionId = sesionId.String()
}
