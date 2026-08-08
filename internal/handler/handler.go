package handler

import (
	"car-valuation-agent/internal/infrastructure/util"
	"car-valuation-agent/internal/model"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/genai"
)

type Handler struct {
	activeAgentRunner *runner.Runner
}

func New(agentRunner runner.Runner) Handler {
	return Handler{
		activeAgentRunner: &agentRunner,
	}
}

func (h Handler) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /chat", h.askToAgent)
	mux.HandleFunc("GET /health", h.healthCheck)
	return mux
}

func (h Handler) askToAgent(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var chatRequest model.ChatRequest
	if err := chatRequest.FromJSON(requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("X-USER-ID") != "" {
		userId, err := uuid.Parse(r.Header.Get("x-user-id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		chatRequest.SetUserId(userId)
	}

	if chatRequest.SessionId == "" {
		sessionId := util.NewSessionId()
		chatRequest.SetSessionId(sessionId)
	}

	if err := chatRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	content := genai.NewContentFromText(chatRequest.Message, genai.RoleUser)

	var agentResponse string
	for event, err := range h.activeAgentRunner.Run(ctx, chatRequest.UserId, chatRequest.SessionId, content, agent.RunConfig{}) {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if event.Content != nil {
			for _, part := range event.Content.Parts {
				if part.Text != "" && !part.Thought {
					agentResponse += part.Text
				}
			}
		}
	}

	var carValuationResponse model.CarValuationResponse
	if err := util.FromJSON(agentResponse, &carValuationResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := util.WriteJSON(w, carValuationResponse, http.StatusOK); err != nil {
		log.Printf("askToAgent[error writing response: %v]", err)
	}
}

func (h Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	var response model.HealthCheckResponse
	status := http.StatusOK
	if h.isAgentRunning() {
		response = model.OfResponse("UP")
	} else {
		response = model.OfResponse("DOWN")
		status = http.StatusServiceUnavailable
	}

	if err := util.WriteJSON(w, response, status); err != nil {
		log.Printf("healthCheck[error writing response: %v]", err)
	}
}

func (h Handler) isAgentRunning() bool {
	return h.activeAgentRunner != nil
}
