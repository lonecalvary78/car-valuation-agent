package config

import (
	"os"
	"strconv"
	"time"
)

func Load() AppConfig {
	return AppConfig{
		server: Server{
			Host:        os.Getenv("SERVER_HOST"),
			Port:        getEnvAsInt(os.Getenv("SERVER_PORT"), 8080),
			ReadTimeout: getEnvAsDuration(os.Getenv("SERVER_READ_TIMEOUT"), 30*time.Second),
			Writetimout: getEnvAsDuration(os.Getenv("SERVER_WRITE_TIMEOUT"), 2*time.Minute),
		},
		agent: Agent{
			Name: os.Getenv("AGENT_NAME"),
			model: Model{
				BaseUrl:   os.Getenv("BASE_URL"),
				ApiKey:    os.Getenv("API_KEY"),
				ModelName: os.Getenv("MODEL_NAME"),
			},
			skill: Skill{
				Name:     os.Getenv("SKILL_NAME"),
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

func getEnvAsDuration(envVariableValue string, defaultValue time.Duration) time.Duration {
	convertedValue, err := time.ParseDuration(envVariableValue)
	if err != nil {
		return defaultValue
	}
	return convertedValue
}
