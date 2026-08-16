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
	validationErrors := make([]error, 0)
	var validationError error
	//Server Config
	if validationError = validator.ValidateForRequiredOfString("Host", config.GetServer().Host); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}
	if validationError = validator.ValidateForRequiredOfNumeric("Port", config.GetServer().Port); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}
	if validationError = validator.ValidateForRequiredOfDuration("Read TimeOut", config.GetServer().ReadTimeout); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}
	if validationError = validator.ValidateForRequiredOfDuration("Write TimeOut", config.GetServer().WriteTimeout); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	//Agent Config
	if validationError = validator.ValidateForRequiredOfString("Agent Name", config.GetAgent().Name); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	//Model Config
	if validationError = validator.ValidateForRequiredOfString("BaseUrl", config.GetAgent().GetModel().BaseUrl); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	if validationError = validator.ValidateForRequiredOfString("ApiKey", config.GetAgent().GetModel().ApiKey); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	if validationError = validator.ValidateForRequiredOfString("Model Name", config.GetAgent().GetModel().ModelName); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	if validationError = validator.ValidateForRequiredOfString("Skill - Location", config.GetAgent().GetSkill().Location); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	//Keycloak Config
	if validationError = validator.ValidateForRequiredOfString("Keycloak - BaseUrl", config.GetKeycloak().BaseUrl); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	if validationError = validator.ValidateForRequiredOfString("Keycloak - Realm", config.GetKeycloak().Realm); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	if validationError = validator.ValidateForRequiredOfString("Keycloak - ClientId", config.GetKeycloak().ClientId); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	if validationError = validator.ValidateForRequiredOfDuration("Wait Timeout", config.GetWaitTimeout()); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	//Redis Config
	if validationError = validator.ValidateForRequiredOfString("Redis - Addr", config.GetRedis().Addr); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	//RateLimit Config
	if validationError = validator.ValidateForRequiredOfNumeric("RateLimit - Limit", config.GetRateLimit().Limit); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	if validationError = validator.ValidateForRequiredOfDuration("RateLimit - Window", config.GetRateLimit().Window); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	return errors.Join(validationErrors...)
}
