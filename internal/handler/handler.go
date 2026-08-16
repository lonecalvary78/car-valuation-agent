package handler

import (
	"car-valuation-agent/internal/infrastructure/middleware"
	"car-valuation-agent/internal/infrastructure/util"
	"car-valuation-agent/internal/model"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/genai"
)

type Handler struct {
	activeAgentRunner *runner.Runner
	waitTimeout       time.Duration
}

func New(agentRunner runner.Runner, waitTimeout time.Duration) Handler {
	return Handler{
		activeAgentRunner: &agentRunner,
		waitTimeout:       waitTimeout,
	}
}

func (h Handler) RegisterRoutes(authMiddleware middleware.Middleware) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("POST /v1/valuations", authMiddleware(http.HandlerFunc(h.askToAgent)))
	mux.HandleFunc("GET /health", h.healthCheck)
	return mux
}

func (h Handler) askToAgent(w http.ResponseWriter, r *http.Request) {
	chatRequest, err := parseChatRequest(w, r)
	if err != nil {
		return
	}

	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "authenticated user not found in request context", http.StatusUnauthorized)
		return
	}
	chatRequest.UserId = user.ID

	ctx, cancel := context.WithTimeout(r.Context(), h.waitTimeout)
	defer cancel()

	agentResponse, err := h.runAgent(ctx, chatRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var carValuationResponse model.CarValuationResponse
	err = util.FromJSON(util.ExtractJSONObject(agentResponse), &carValuationResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = util.WriteJSON(w, carValuationResponse, http.StatusOK)
	if err != nil {
		log.Printf("askToAgent[error writing response: %v]", err)
	}
}

// parseChatRequest reads and validates the incoming request, writing an error response and
// returning a non-nil error when the request is malformed or invalid.
func parseChatRequest(w http.ResponseWriter, r *http.Request) (model.ChatRequest, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 32<<10)
	bodyContents, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return model.ChatRequest{}, fmt.Errorf("handler: failed to read request body: %w", err)
	}

	var chatRequest model.ChatRequest
	err = chatRequest.FromJSON(bodyContents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return model.ChatRequest{}, fmt.Errorf("handler: failed to parse request body: %w", err)
	}

	if chatRequest.SessionId == "" {
		chatRequest.SetSessionId(util.NewSessionId())
	}

	err = chatRequest.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return model.ChatRequest{}, fmt.Errorf("handler: invalid chat request: %w", err)
	}

	return chatRequest, nil
}

func (h Handler) runAgent(ctx context.Context, chatRequest model.ChatRequest) (string, error) {
	content := genai.NewContentFromText(chatRequest.Message, genai.RoleUser)

	var agentResponse strings.Builder
	for event, err := range h.activeAgentRunner.Run(ctx, chatRequest.UserId, chatRequest.SessionId, content, agent.RunConfig{}) {
		if err != nil {
			return "", err
		}

		if event.Content == nil {
			continue
		}

		for _, part := range event.Content.Parts {
			if part.Text != "" && !part.Thought {
				agentResponse.WriteString(part.Text)
			}
		}
	}

	return agentResponse.String(), nil
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

	err := util.WriteJSON(w, response, status)
	if err != nil {
		log.Printf("healthCheck[error writing response: %v]", err)
	}
}

func (h Handler) isAgentRunning() bool {
	return h.activeAgentRunner != nil
}
