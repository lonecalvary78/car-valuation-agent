package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimit(client *redis.Client, limit int, window time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := "ratelimit:" + rateLimitIdentifier(r)

			count, ttl, err := incrementCounter(r.Context(), client, key, window)
			if err != nil {
				log.Printf("RateLimit[error: %v]", err)
				next.ServeHTTP(w, r)
				return
			}

			remaining := limit - int(count)
			if remaining < 0 {
				remaining = 0
			}
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			if count > int64(limit) {
				w.Header().Set("Retry-After", strconv.Itoa(int(ttl.Seconds())))
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func incrementCounter(ctx context.Context, client *redis.Client, key string, window time.Duration) (int64, time.Duration, error) {
	pipe := client.TxPipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, 0, err
	}

	ttl, err := client.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		ttl = window
	}

	return incrCmd.Val(), ttl, nil
}

func rateLimitIdentifier(r *http.Request) string {
	if user, ok := r.Context().Value(UserKey).(User); ok && user.ID != "" {
		return "user:" + user.ID
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	return "ip:" + host
}
