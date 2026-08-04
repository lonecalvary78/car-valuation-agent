package model

import (
	"car-valuation-agent/internal/infrastructure/validator"
	"errors"

	"github.com/google/uuid"

	"github.com/bytedance/sonic"
)

type ChatRequest struct {
	SessionId string `json:"sessionId"`
	UserId    string `json:"userId"`
	Message   string ` json:"message"`
}

func (chat *ChatRequest) FromJSON(requestBody []byte) error {
	return sonic.Unmarshal(requestBody, chat)
}

func (chat *ChatRequest) Validate() error {
	validationErrors := make([]error, 0)
	if validationError := validator.ValidateForRequiredOfString("UserId", chat.UserId); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}
	if validationError := validator.ValidateForRequiredOfString("SessionId", chat.SessionId); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	return errors.Join(validationErrors...)
}

func (chat *ChatRequest) SetSesionId(sesionId uuid.UUID) {
	chat.SessionId = sesionId.String()
}

func (chat *ChatRequest) SetUserId(userId uuid.UUID) {
	chat.UserId = userId.String()
}
