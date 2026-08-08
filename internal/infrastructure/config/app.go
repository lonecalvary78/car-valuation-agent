package config

import (
	"car-valuation-agent/internal/infrastructure/validator"
	"errors"
)

type AppConfig struct {
	server Server
	agent  Agent
	redis  Redis
}

func (config AppConfig) GetServer() Server {
	return config.server
}

func (config AppConfig) GetAgent() Agent {
	return config.agent
}

func (config AppConfig) GetRedis() Redis {
	return config.redis
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
	if validationError = validator.ValidateForRequiredOfDuration("Write TimeOut", config.GetServer().Writetimout); validationError != nil {
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

	if validationError = validator.ValidateForRequiredOfString("Model Name", config.GetAgent().GetModel().ApiKey); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	if validationError = validator.ValidateForRequiredOfString("Skill - Location", config.GetAgent().GetSkill().Location); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	//Redis Config
	if validationError = validator.ValidateForRequiredOfString("Host", config.GetRedis().Host); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}
	if validationError = validator.ValidateForRequiredOfNumeric("Port", config.GetRedis().Port); validationError != nil {
		validationErrors = append(validationErrors, validationError)
	}

	return errors.Join(validationErrors...)
}
