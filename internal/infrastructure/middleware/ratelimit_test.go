package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRateLimit_AllowsUnderLimit(t *testing.T) {
	ctx := context.Background()
	client, terminate := createRedisClient(ctx, t)
	defer terminate()

	handler := RateLimit(client, 2, time.Minute)(okHandler())

	for range 2 {
		recorder := performRequest(t, handler, "203.0.113.10:1234")
		require.Equal(t, http.StatusOK, recorder.Code)
	}
}

func TestRateLimit_RejectsOverLimit(t *testing.T) {
	ctx := context.Background()
	client, terminate := createRedisClient(ctx, t)
	defer terminate()

	handler := RateLimit(client, 2, time.Minute)(okHandler())

	for range 2 {
		recorder := performRequest(t, handler, "203.0.113.20:1234")
		require.Equal(t, http.StatusOK, recorder.Code)
	}

	recorder := performRequest(t, handler, "203.0.113.20:1234")
	require.Equal(t, http.StatusTooManyRequests, recorder.Code)
	require.NotEmpty(t, recorder.Header().Get("Retry-After"))
}

func TestRateLimit_TracksCallersIndependently(t *testing.T) {
	ctx := context.Background()
	client, terminate := createRedisClient(ctx, t)
	defer terminate()

	handler := RateLimit(client, 1, time.Minute)(okHandler())

	require.Equal(t, http.StatusOK, performRequest(t, handler, "203.0.113.30:1234").Code)
	require.Equal(t, http.StatusTooManyRequests, performRequest(t, handler, "203.0.113.30:1234").Code)
	require.Equal(t, http.StatusOK, performRequest(t, handler, "203.0.113.31:1234").Code)
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func performRequest(t *testing.T, handler http.Handler, remoteAddr string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/v1/valuations", nil)
	request.RemoteAddr = remoteAddr
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func createRedisClient(ctx context.Context, t *testing.T) (*redis.Client, func()) {
	t.Helper()
	container, err := rediscontainer.Run(ctx, "redis:alpine")
	require.NoError(t, err)

	connectionString, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	parsedAddr, err := url.Parse(connectionString)
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{Addr: parsedAddr.Host})

	return client, func() {
		closeErr := client.Close()
		if closeErr != nil {
			t.Errorf("failed to close redis client: %v", closeErr)
		}
		terminateErr := container.Terminate(ctx)
		if terminateErr != nil {
			t.Errorf("failed to terminate redis container: %v", terminateErr)
		}
	}
}
