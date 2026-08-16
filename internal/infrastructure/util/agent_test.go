package util

import (
	"car-valuation-agent/internal/apptest"
	"car-valuation-agent/internal/infrastructure/config"
	"context"
	"fmt"
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

	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("agent_test: mocked tools call failed: %w", err)
	}

	return tools, nil
}

func TestOfAgent(t *testing.T) {
	mockTools := new(AgentToolMock)
	apptest.SetEnvForTesting(t)
	appConfig, err := config.Load()
	require.NoError(t, err)
	require.NoError(t, err)
	createdAgent, err := OfAgent(context.Background(), "Do something", appConfig, []tool.Toolset{mockTools})
	require.NoError(t, err)
	require.Equal(t, "car-valuation-agent", createdAgent.Name())
}
