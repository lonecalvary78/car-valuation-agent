package config

import (
	"os"
	"strconv"
	"time"
)

func Load() (AppConfig, error) {
	readTimeOut, writeTimeout, waitTimeout, err := loadAllTimeouts()
	if err != nil {
		return AppConfig{}, err
	}
	appConfig := AppConfig{
		server: Server{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvAsInt(os.Getenv("SERVER_PORT"), 8080),
			ReadTimeout:  readTimeOut,
			WriteTimeout: writeTimeout,
		},
		agent: Agent{
			Name: os.Getenv("AGENT_NAME"),
			model: Model{
				BaseUrl:   os.Getenv("BASE_URL"),
				ApiKey:    os.Getenv("API_KEY"),
				ModelName: os.Getenv("MODEL_NAME"),
				AdvanceSetup: AdvancedSetup{
					Temperature:      getEnvAsFloat(os.Getenv("TEMPERATURE"), 0.0),
					TopP:             getEnvAsFloat(os.Getenv("TOP_P"), 0.0),
					TopK:             getEnvAsFloat(os.Getenv("TOP_K"), 0.0),
					PresencePenalty:  getEnvAsFloat(os.Getenv("PRESENCE_PENALTY"), 0.0),
					FrequencyPenalty: getEnvAsFloat(os.Getenv("FREQUENCY_PENALTY"), 0.0),
					MaximumTokens:    int32(getEnvAsInt(os.Getenv("MAXIMUM_TOKENS"), 0)),
				},
			},
			skill: Skill{
				Location: os.Getenv("SKILL_LOCATION"),
			},
		},
		waitTimeout: waitTimeout,
	}
	if err := appConfig.Validate(); err != nil {
		return AppConfig{}, err
	}
	return appConfig, nil
}

func loadAllTimeouts() (time.Duration, time.Duration, time.Duration, error) {
	readTimeout, err := getEnvAsDuration(os.Getenv("SERVER_READ_TIMEOUT"), 30*time.Second)
	if err != nil {
		return 0 * time.Second, 0 * time.Second, 0 * time.Second, err
	}

	writeTimeout, err := getEnvAsDuration(os.Getenv("SERVER_WRITE_TIMEOUT"), 2*time.Minute)
	if err != nil {
		return 0 * time.Second, 0 * time.Second, 0 * time.Second, err
	}

	waitTimeout, err := getEnvAsDuration(os.Getenv("WAIT_TIMEOUT"), 60*time.Second)
	if err != nil {
		return 0 * time.Second, 0 * time.Second, 0 * time.Second, err
	}

	return readTimeout, writeTimeout, waitTimeout, nil
}

func getEnv(envVariableName string, defaultValue string) string {
	envVariableValue := os.Getenv(envVariableName)
	if envVariableValue == "" {
		return defaultValue
	}
	return envVariableValue
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

func getEnvAsDuration(envVariableValue string, defaultValue time.Duration) (time.Duration, error) {
	if envVariableValue == "" {
		return defaultValue, nil
	}

	convertedValue, err := time.ParseDuration(envVariableValue)
	if err != nil {
		return 0 * time.Second, err
	} else if convertedValue < 0*time.Second {
		return defaultValue, nil
	}
	return convertedValue, nil
}
