package config

import (
	"car-valuation-agent/internal/apptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	apptest.SetEnvForTesting(t)
	appConfig, err := Load()
	require.NoError(t, err)
	require.Equal(t, os.Getenv("SERVER_HOST"), appConfig.GetServer().Host)
	require.Equal(t, getEnvAsInt(os.Getenv("SERVER_PORT"), 8585), appConfig.GetServer().Port)
	require.Equal(t, os.Getenv("AGENT_NAME"), appConfig.GetAgent().Name)
	require.Equal(t, os.Getenv("BASE_URL"), appConfig.GetAgent().GetModel().BaseUrl)
	require.Equal(t, os.Getenv("API_KEY"), appConfig.GetAgent().GetModel().ApiKey)
	require.Equal(t, os.Getenv("MODEL_NAME"), appConfig.GetAgent().GetModel().ModelName)
	require.InDelta(t, getEnvAsFloat(os.Getenv("TEMPERATURE"), 2.0), appConfig.GetAgent().GetModel().AdvanceSetup.Temperature, 0.0001)
	require.InDelta(t, getEnvAsFloat(os.Getenv("TOP_K"), 0.0), appConfig.GetAgent().GetModel().AdvanceSetup.TopK, 0.0001)
	require.InDelta(t, getEnvAsFloat(os.Getenv("TOP_P"), 0.0), appConfig.GetAgent().GetModel().AdvanceSetup.TopP, 0.0001)
	require.InDelta(t, getEnvAsFloat(os.Getenv("PRESENCE_PENALTY"), 1.0), appConfig.GetAgent().GetModel().AdvanceSetup.PresencePenalty, 0.0001)
	require.InDelta(t, getEnvAsFloat(os.Getenv("FREQUENCY_PENALTY"), 1.0), appConfig.GetAgent().GetModel().AdvanceSetup.FrequencyPenalty, 0.0001)
	require.Equal(t, os.Getenv("SKILL_LOCATION"), appConfig.GetAgent().GetSkill().Location)
	require.Equal(t, os.Getenv("KEYCLOAK_BASE_URL"), appConfig.GetKeycloak().BaseUrl)
	require.Equal(t, os.Getenv("KEYCLOAK_REALM"), appConfig.GetKeycloak().Realm)
	require.Equal(t, os.Getenv("KEYCLOAK_CLIENT_ID"), appConfig.GetKeycloak().ClientId)
	require.Equal(t, os.Getenv("REDIS_ADDR"), appConfig.GetRedis().Addr)
	require.Equal(t, os.Getenv("REDIS_PASSWORD"), appConfig.GetRedis().Password)
	require.Equal(t, getEnvAsInt(os.Getenv("REDIS_DB"), -1), appConfig.GetRedis().DB)
	require.Equal(t, getEnvAsInt(os.Getenv("RATE_LIMIT_REQUESTS"), -1), appConfig.GetRateLimit().Limit)
	expectedRateLimitWindow, err := time.ParseDuration(os.Getenv("RATE_LIMIT_WINDOW"))
	require.NoError(t, err)
	require.Equal(t, expectedRateLimitWindow, appConfig.GetRateLimit().Window)
}
