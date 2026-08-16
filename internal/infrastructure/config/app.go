package config

import (
	"car-valuation-agent/internal/infrastructure/validator"
	"errors"
	"time"
)

type AppConfig struct {
	server      Server
	agent       Agent
	keycloak    Keycloak
	redis       Redis
	rateLimit   RateLimit
	waitTimeout time.Duration
}

func (config AppConfig) GetServer() Server {
	return config.server
}

func (config AppConfig) GetAgent() Agent {
	return config.agent
}

func (config AppConfig) GetKeycloak() Keycloak {
	return config.keycloak
}

func (config AppConfig) GetRedis() Redis {
	return config.redis
}

func (config AppConfig) GetRateLimit() RateLimit {
	return config.rateLimit
}

func (config AppConfig) GetWaitTimeout() time.Duration {
	return config.waitTimeout
}

func (config AppConfig) Validate() error {
	validationErrors := []error{
		// Server Config
		validator.ValidateForRequiredOfString("Host", config.GetServer().Host),
		validator.ValidateForRequiredOfNumeric("Port", config.GetServer().Port),
		validator.ValidateForRequiredOfDuration("Read TimeOut", config.GetServer().ReadTimeout),
		validator.ValidateForRequiredOfDuration("Write TimeOut", config.GetServer().WriteTimeout),

		// Agent Config
		validator.ValidateForRequiredOfString("Agent Name", config.GetAgent().Name),

		// Model Config
		validator.ValidateForRequiredOfString("BaseUrl", config.GetAgent().GetModel().BaseUrl),
		validator.ValidateForRequiredOfString("ApiKey", config.GetAgent().GetModel().ApiKey),
		validator.ValidateForRequiredOfString("Model Name", config.GetAgent().GetModel().ModelName),
		validator.ValidateForRequiredOfString("Skill - Location", config.GetAgent().GetSkill().Location),

		// Keycloak Config
		validator.ValidateForRequiredOfString("Keycloak - BaseUrl", config.GetKeycloak().BaseUrl),
		validator.ValidateForRequiredOfString("Keycloak - Realm", config.GetKeycloak().Realm),
		validator.ValidateForRequiredOfString("Keycloak - ClientId", config.GetKeycloak().ClientId),
		validator.ValidateForRequiredOfDuration("Wait Timeout", config.GetWaitTimeout()),

		// Redis Config
		validator.ValidateForRequiredOfString("Redis - Addr", config.GetRedis().Addr),

		// RateLimit Config
		validator.ValidateForRequiredOfNumeric("RateLimit - Limit", config.GetRateLimit().Limit),
		validator.ValidateForRequiredOfDuration("RateLimit - Window", config.GetRateLimit().Window),
	}

	return errors.Join(validationErrors...)
}
