package sessionservice

import (
	"fmt"

	redissession "github.com/achetronic/adk-utils-go/session/redis"
)

// OfRedisBackedSessionService creates a Redis backed SessionService.
func OfRedisBackedSessionService(redisAddr string, redisDB int) (*redissession.RedisSessionService, error) {
	sessionService, err := redissession.NewRedisSessionService(redissession.RedisSessionServiceConfig{
		Addr: redisAddr,
		DB:   redisDB,
	})
	if err != nil {
		return nil, fmt.Errorf("sessionservice: failed to create redis session service: %w", err)
	}

	return sessionService, nil
}
