package model

import "github.com/bytedance/sonic"

type ChatRequest struct {
	SessionId string `json:"sessionId"`
	UserId    string `json:"userId"`
	Message   string ` json:"message"`
}

func (chat *ChatRequest) FromJSON(requestBody []byte) error {
	return sonic.Unmarshal(requestBody, chat)
}
