package handler

import (
	"car-valuation-agent/internal/infrastructure/util"
	"car-valuation-agent/internal/model"
	"io"
	"log"
	"net/http"

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

	if r.Header.Get("x-user-id") != "" {
		userId, _ := uuid.Parse(r.Header.Get("x-user-id"))
		chatRequest.SetUserId(userId)
	}

	if chatRequest.SessionId == "" {
		sessionId := util.NewSessionId()
		chatRequest.SetSesionId(sessionId)
	}

	content := genai.NewContentFromText(chatRequest.Message, genai.RoleUser)
	var response string
	for event, err := range h.activeAgentRunner.Run(r.Context(), chatRequest.UserId, chatRequest.SessionId, content, agent.RunConfig{}) {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if event.Content != nil {
			for _, part := range event.Content.Parts {
				if part.Text != "" && !part.Thought {
					response += part.Text
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(response)); err != nil {
		log.Printf("askToAgent[error writing response: %v]", err)
	}
}

func (h Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	var response model.HealthCheckResponse
	status := http.StatusOK
	if h.IsAagentRunning() {
		response = model.OfResponse("UP")
	} else {
		response = model.OfResponse("DOWN")
		status = http.StatusServiceUnavailable
	}

	responseBytes, err := response.ToJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(responseBytes); err != nil {
		log.Printf("healthCheck[error writing response: %v]", err)
	}
}

func (h Handler) IsAagentRunning() bool {
	return h.activeAgentRunner != nil
}
