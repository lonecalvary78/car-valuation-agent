package handler

import (
	"car-valuation-agent/internal/infrastructure/util"
	"car-valuation-agent/internal/model"
	"io"
	"net/http"

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

	if chatRequest.SessionId == "" {
		sessionId := util.NewSessionId()
		chatRequest.SessionId = sessionId.String()
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
	w.Write([]byte(response))

}

func (h Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	if h.activeAgentRunner != nil {
		response := model.OfResponse("OK")
		responseBytes, err := response.ToJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write(responseBytes)
	}
}
