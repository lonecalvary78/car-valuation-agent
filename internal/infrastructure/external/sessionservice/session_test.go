package sessionservice

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"google.golang.org/adk/v2/session"
)

func TestCreateNewSession(t *testing.T) {
	ctx := context.Background()
	redis, redisAddr, err := createRedisContainer(ctx)
	require.NoError(t, err)

	defer func() { require.NoError(t, redis.Terminate(ctx)) }()

	redisBacedSessionService, err := OfRedisBackedSessionService(redisAddr, 0)
	require.NoError(t, err)

	defer func() { require.NoError(t, redisBacedSessionService.Close()) }()
	response, err := redisBacedSessionService.Create(ctx, &session.CreateRequest{
		SessionID: "4767576",
		UserID:    "12345",
		AppName:   "test-app",
	})
	require.NoError(t, err)
	require.Equal(t, "4767576", response.Session.ID())
	require.Equal(t, "test-app", response.Session.AppName())
}

func createRedisContainer(ctx context.Context) (*rediscontainer.RedisContainer, string, error) {
	redis, err := rediscontainer.Run(ctx, "redis:alpine", rediscontainer.WithLogLevel(rediscontainer.LogLevelDebug))
	if err != nil {
		return nil, "", fmt.Errorf("session_test: failed to start redis container: %w", err)
	}

	redisConnStr, err := redis.ConnectionString(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("session_test: failed to get connection string: %w", err)
	}

	parsedAddr, err := url.Parse(redisConnStr)
	if err != nil {
		return nil, "", fmt.Errorf("session_test: failed to parse connection string: %w", err)
	}

	return redis, parsedAddr.Host, nil
}
