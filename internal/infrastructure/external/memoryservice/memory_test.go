package memoryservice

import (
	"car-valuation-agent/internal/infrastructure/config"
	"car-valuation-agent/internal/infrastructure/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/achetronic/adk-utils-go/memory/postgres"
	"github.com/stretchr/testify/require"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"google.golang.org/adk/v2/memory"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/session"
	"google.golang.org/genai"
)

const embeddingModelName = "any-model"

func newMockEmbeddingServer(t *testing.T) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		embedding := make([]float32, 8)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"embedding": embedding, "index": 0},
			},
			"model": embeddingModelName,
		})
		if err != nil {
			t.Errorf("failed to encode mock embedding response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	return server
}

func TestAddMemory(t *testing.T) {
	ctx := context.Background()
	postgresContainer, dbConnStr, err := createPostgresContrainer(ctx)
	require.NoError(t, err)
	defer func() { require.NoError(t, postgresContainer.Terminate(ctx)) }()

	embeddingServer := newMockEmbeddingServer(t)

	pgBackedMemoryService, err := OfExternalMemoryService(ctx, dbConnStr, config.Model{
		BaseUrl:   embeddingServer.URL,
		ApiKey:    "dummyKey",
		ModelName: embeddingModelName,
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, pgBackedMemoryService.Close()) }()

	sessionService := session.InMemoryService()
	createResp, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName:   "test-app",
		UserID:    "user-1",
		SessionID: "sess-1",
	})
	require.NoError(t, err)

	err = sessionService.AppendEvent(ctx, createResp.Session, &session.Event{
		ID:     "event-1",
		Author: "user",
		LLMResponse: model.LLMResponse{
			Content: &genai.Content{
				Role:  "user",
				Parts: []*genai.Part{genai.NewPartFromText("What is the capital of France?")},
			},
		},
	})
	require.NoError(t, err)

	err = pgBackedMemoryService.AddSessionToMemory(ctx, createResp.Session)
	require.NoError(t, err)
}

func TestSearchMemory(t *testing.T) {
	ctx := context.Background()
	postgresContainer, dbConnStr, err := createPostgresContrainer(ctx)
	require.NoError(t, err)
	defer func() { require.NoError(t, postgresContainer.Terminate(ctx)) }()

	embeddingServer := newMockEmbeddingServer(t)

	pgBackedMemoryService, err := OfExternalMemoryService(ctx, dbConnStr, config.Model{
		BaseUrl:   embeddingServer.URL,
		ApiKey:    "dummyKey",
		ModelName: embeddingModelName,
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, pgBackedMemoryService.Close()) }()

	require.NoError(t, prepareMemoryForTest(ctx, pgBackedMemoryService, "test-app", "1232"))

	searchResponse, err := pgBackedMemoryService.SearchMemory(ctx, &memory.SearchRequest{
		UserID:  "1232",
		AppName: "test-app",
		Query:   "",
	})
	require.NoError(t, err)
	require.NotEmpty(t, searchResponse.Memories)
}

func prepareMemoryForTest(ctx context.Context, pgBackedMemoryService *postgres.PostgresMemoryService, appName string, userId string) error {
	sessionService := session.InMemoryService()
	createResp, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName:   appName,
		UserID:    userId,
		SessionID: util.NewSessionId().String(),
	})
	if err != nil {
		return fmt.Errorf("memory_test: failed to create session: %w", err)
	}

	err = sessionService.AppendEvent(ctx, createResp.Session, &session.Event{
		ID:     "event-1",
		Author: "user",
		LLMResponse: model.LLMResponse{
			Content: &genai.Content{
				Role:  "user",
				Parts: []*genai.Part{genai.NewPartFromText("What is the capital of France?")},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("memory_test: failed to append event: %w", err)
	}

	err = pgBackedMemoryService.AddSessionToMemory(ctx, createResp.Session)
	if err != nil {
		return fmt.Errorf("memory_test: failed to add session to memory: %w", err)
	}

	return nil
}

func createPostgresContrainer(ctx context.Context) (*pgcontainer.PostgresContainer, string, error) {
	pgContainer, err := pgcontainer.Run(ctx, "pgvector/pgvector:pg18", pgcontainer.WithDatabase("memory-db"),
		pgcontainer.WithUsername("postgres"), pgcontainer.WithPassword("postgres"), pgcontainer.BasicWaitStrategies())
	if err != nil {
		return nil, "", fmt.Errorf("memory_test: failed to start postgres container: %w", err)
	}

	dbConnStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, "", fmt.Errorf("memory_test: failed to get connection string: %w", err)
	}

	return pgContainer, dbConnStr, nil
}
