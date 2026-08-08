package config

import "testing"

func SetEnv(t *testing.T) {
	t.Helper()
	t.Setenv("SERVER_HOST", "localhost")
	t.Setenv("SERVER_PORT", "8585")
	t.Setenv("BASE_URL", "http://localhost:8000")
	t.Setenv("API_KEY", "apikey")
	t.Setenv("AGENT_NAME", "car-valuation-agent")
	t.Setenv("MODEL_NAME", "qwen/qwen3.6_27b")
	t.Setenv("TEMPERATURE", "0.6")
	t.Setenv("TOP_K", "20")
	t.Setenv("TOP_P", "0.95")
	t.Setenv("PRESEMCE_PENALTY", "1.0")
	t.Setenv("FREQENCY_PENALTY", "1.0")
	t.Setenv("SKILL_LOCATION", "skills")
}
