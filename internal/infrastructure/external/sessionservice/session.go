package sessionservice

import (
	redissession "github.com/achetronic/adk-utils-go/session/redis"
)

// Create Redis backed SessionService
func OfRedisBackedSessionService(redisAddr string, redisDB int) (*redissession.RedisSessionService, error) {
	return redissession.NewRedisSessionService(redissession.RedisSessionServiceConfig{
		Addr: redisAddr,
		DB:   redisDB,
	})
}
