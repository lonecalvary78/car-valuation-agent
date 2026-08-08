package util

import (
	"car-valuation-agent/internal/infrastructure/config"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
)

type AgentToolMock struct {
	mock.Mock
}

func (mockedTool *AgentToolMock) Name() string {
	return mockedTool.Called().String(0)
}

func (mockedTool *AgentToolMock) Description() string {
	return mockedTool.Called().String(1)
}

func (mockedTool *AgentToolMock) IsLongRunning() bool {
	return mockedTool.Called().Bool(3)
}

func (mockedTool *AgentToolMock) Tools(ctx agent.ReadonlyContext) ([]tool.Tool, error) {
	args := mockedTool.Called(ctx)
	tools, _ := args.Get(0).([]tool.Tool)
	return tools, args.Error(1)
}

func setEnvForTesting(t *testing.T, targetModel string) config.AppConfig {
	t.Helper()
	t.Setenv("AGENT_NAME", "car-valuer-agent")
	t.Setenv("MODEL_NAME", targetModel)
	t.Setenv("BASE_URL", "http:localhost:8000")
	t.Setenv("API_KEY", "apiKey")
	return config.Load()
}

func TestOfAgent(t *testing.T) {
	mockTools := new(AgentToolMock)
	appConfig := setEnvForTesting(t, "mdoel-123")

	createdAgent, err := OfAgent(context.Background(), "Do something", appConfig, []tool.Toolset{mockTools})
	require.NoError(t, err)
	require.Equal(t, "car-valuer-agent", createdAgent.Name())
}
