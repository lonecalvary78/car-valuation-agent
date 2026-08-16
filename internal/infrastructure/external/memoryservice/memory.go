package memoryservice

import (
	"car-valuation-agent/internal/infrastructure/config"
	"context"
	"fmt"

	pgmemory "github.com/achetronic/adk-utils-go/memory/postgres"
)

// OfExternalMemoryService creates the PostgreSQL backed Memory Service.
func OfExternalMemoryService(ctx context.Context, dbConnStr string, modelConfig config.Model) (*pgmemory.PostgresMemoryService, error) {
	memoryService, err := pgmemory.NewPostgresMemoryService(ctx, pgmemory.PostgresMemoryServiceConfig{
		ConnString: dbConnStr,
		EmbeddingModel: pgmemory.NewOpenAICompatibleEmbedding(pgmemory.OpenAICompatibleEmbeddingConfig{
			BaseURL: modelConfig.BaseUrl,
			APIKey:  modelConfig.ApiKey,
			Model:   modelConfig.ModelName,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("memoryservice: failed to create postgres memory service: %w", err)
	}

	return memoryService, nil
}
