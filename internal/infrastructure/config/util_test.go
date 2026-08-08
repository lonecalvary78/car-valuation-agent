package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	SetEnv(t)
	appConfig := Load()
	require.Equal(t, os.Getenv("SERVER_HOST"), appConfig.GetServer().Host)
	require.Equal(t, getEnvAsInt(os.Getenv("SERVER_PORT"), 8585), appConfig.GetServer().Port)
	require.Equal(t, os.Getenv("AGENT_NAME"), appConfig.GetAgent().Name)
	require.Equal(t, os.Getenv("BASE_URL"), appConfig.GetAgent().GetModel().BaseUrl)
	require.Equal(t, os.Getenv("API_KEY"), appConfig.GetAgent().GetModel().ApiKey)
	require.Equal(t, os.Getenv("MODEL_NAME"), appConfig.GetAgent().GetModel().ModelName)
	require.Equal(t, getEnvAsFloat(os.Getenv("TEMPERATURE"), 2.0), appConfig.GetAgent().GetModel().AdvanceSetup.Temperature)
	require.Equal(t, getEnvAsFloat(os.Getenv("TOP_K"), 0.0), appConfig.GetAgent().GetModel().AdvanceSetup.TopK)
	require.Equal(t, getEnvAsFloat(os.Getenv("TOP_P"), 0.0), appConfig.GetAgent().GetModel().AdvanceSetup.TopP)
	require.Equal(t, getEnvAsFloat(os.Getenv("PRESENCE_PENALTY"), 0.0), appConfig.GetAgent().GetModel().AdvanceSetup.PresencePenalty)
	require.Equal(t, getEnvAsFloat(os.Getenv("FREQUENCY_PENALTY"), 0.0), appConfig.GetAgent().GetModel().AdvanceSetup.FrequencyPenalty)
	require.Equal(t, os.Getenv("SKILL_LOCATION"), appConfig.GetAgent().GetSkill().Location)
}
