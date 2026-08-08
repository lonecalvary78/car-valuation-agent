package util

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/memory"
	"google.golang.org/adk/v2/session"
)

func newAgent(t *testing.T) agent.Agent {
	t.Helper()
	createdAgent, _ := agent.New(agent.Config{
		Name:        "car-valuation-agent",
		Description: "car-valuation-agent",
		Run: func(ic agent.InvocationContext) iter.Seq2[*session.Event, error] {
			return func(yield func(*session.Event, error) bool) {}
		},
	})
	return createdAgent
}

func TestCreateRunner(t *testing.T) {
	agent := newAgent(t)
	runner, err := OfRunner(agent, memory.InMemoryService(), session.InMemoryService(), true)
	require.NoError(t, err)
	require.Equal(t, "car-valuation-agent", agent.Name())
	require.NotEmpty(t, runner)
}
