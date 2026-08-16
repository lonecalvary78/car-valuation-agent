package memoryservice

import (
	"car-valuation-agent/internal/infrastructure/config"
	"context"

	pgmemory "github.com/achetronic/adk-utils-go/memory/postgres"
)

// Creating the PostgreSQL backed Memory Service
func OfExternalMemoryService(ctx context.Context, dbConnStr string, modelConfig config.Model) (*pgmemory.PostgresMemoryService, error) {
	return pgmemory.NewPostgresMemoryService(ctx, pgmemory.PostgresMemoryServiceConfig{
		ConnString: dbConnStr,
		EmbeddingModel: pgmemory.NewOpenAICompatibleEmbedding(pgmemory.OpenAICompatibleEmbeddingConfig{
			BaseURL: modelConfig.BaseUrl,
			APIKey:  modelConfig.ApiKey,
			Model:   modelConfig.ModelName,
		}),
	})
}
