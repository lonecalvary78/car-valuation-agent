package config

import (
	"os"
	"strconv"
	"time"
)

func Load() AppConfig {
	return AppConfig{
		server: Server{
			Host:         os.Getenv("SERVER_HOST"),
			Port:         getEnvAsInt(os.Getenv("SERVER_PORT"), 8080),
			ReadTimeout:  getEnvAsDuration(os.Getenv("SERVER_READ_TIMEOUT"), 30*time.Second),
			WriteTimeout: getEnvAsDuration(os.Getenv("SERVER_WRITE_TIMEOUT"), 2*time.Minute),
		},
		agent: Agent{
			Name: os.Getenv("AGENT_NAME"),
			model: Model{
				BaseUrl:   os.Getenv("BASE_URL"),
				ApiKey:    os.Getenv("API_KEY"),
				ModelName: os.Getenv("MODEL_NAME"),
				AdvanceSetup: AdvanceSetup{
					Temperature:      getEnvAsFloat(os.Getenv("TEMPERATURE"), 0.0),
					TopP:             getEnvAsFloat(os.Getenv("TOP_P"), 0.0),
					TopK:             getEnvAsFloat(os.Getenv("TOP_K"), 0.0),
					PresencePenalty:  getEnvAsFloat(os.Getenv("PRESENCE_PENALTY"), 0.0),
					FrequencyPenalty: getEnvAsFloat(os.Getenv("FREQUENCY_PENALTY"), 0.0),
					MaximumTokens:    int32(getEnvAsInt(os.Getenv("MAX_TOKENS"), 0)),
				},
			},
			skill: Skill{
				Location: os.Getenv("SKILL_LOCATION"),
			},
		},
	}
}

func getEnvAsInt(envVariableValue string, defaultValue int) int {
	convertedEnvVariableValue, err := strconv.Atoi(envVariableValue)
	if err != nil {
		return defaultValue
	}
	return convertedEnvVariableValue
}

func getEnvAsFloat(envVariableValue string, defaultValue float32) float32 {
	convertedEnvVariableValue, err := strconv.ParseFloat(envVariableValue, 32)
	if err != nil {
		return defaultValue
	}
	return float32(convertedEnvVariableValue)
}

func getEnvAsDuration(envVariableValue string, defaultValue time.Duration) time.Duration {
	convertedValue, err := time.ParseDuration(envVariableValue)
	if err != nil {
		return defaultValue
	}
	return convertedValue
}
