package util

import (
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/memory"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/session"
)

func OfRunner(targetAgent agent.Agent, memoryService memory.Service, sessionService session.Service, withAutoCreateSession bool) (*runner.Runner, error) {
	return runner.New(runner.Config{
		AppName:           targetAgent.Name(),
		Agent:             targetAgent,
		SessionService:    sessionService,
		MemoryService:     memoryService,
		AutoCreateSession: withAutoCreateSession,
	})
}
