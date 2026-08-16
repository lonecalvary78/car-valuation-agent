package apptest

import "testing"

func SetEnvForTesting(t *testing.T) {
	t.Helper()
	t.Setenv("SERVER_HOST", "localhost")
	t.Setenv("SERVER_PORT", "8585")
	t.Setenv("SERVER_READ_TIMEOUT", "30s")
	t.Setenv("SERVER_WRITE_TIMEOUT", "2m")
	t.Setenv("BASE_URL", "http://localhost:8000")
	t.Setenv("API_KEY", "apikey")
	t.Setenv("AGENT_NAME", "car-valuation-agent")
	t.Setenv("MODEL_NAME", "qwen/qwen3.6_27b")
	t.Setenv("TEMPERATURE", "0.6")
	t.Setenv("TOP_K", "20")
	t.Setenv("TOP_P", "0.95")
	t.Setenv("PRESENCE_PENALTY", "1.0")
	t.Setenv("FREQUENCY_PENALTY", "1.0")
	t.Setenv("MAXIMUM_TOKENS", "2048")
	t.Setenv("SKILL_LOCATION", "skills")
	t.Setenv("WAIT_TIMEOUT", "60s")
	t.Setenv("KEYCLOAK_BASE_URL", "http://localhost:8180")
	t.Setenv("KEYCLOAK_REALM", "car-valuation-agent")
	t.Setenv("KEYCLOAK_CLIENT_ID", "car-valuation-agent")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "0")
	t.Setenv("RATE_LIMIT_REQUESTS", "60")
	t.Setenv("RATE_LIMIT_WINDOW", "1m")
}
