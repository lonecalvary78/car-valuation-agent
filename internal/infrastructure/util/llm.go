package util

import (
	"car-valuation-agent/internal/infrastructure/config"
	"context"
	"os"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/memory"
	"google.golang.org/adk/v2/model/openaimodel"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/skilltoolset"
	"google.golang.org/adk/v2/tool/skilltoolset/skill"
)

func OfAgent(ctx context.Context, instruction string, appConfig config.AppConfig, tools []tool.Toolset) (agent.Agent, error) {
	llm, err := openaimodel.NewModel(ctx, appConfig.GetAgent().GetModel().ModelName, &openaimodel.ClientConfig{
		BaseURL: appConfig.GetAgent().GetModel().BaseUrl,
		APIKey:  appConfig.GetAgent().GetModel().ApiKey,
	})
	if err != nil {
		return nil, err
	}

	return llmagent.New(llmagent.Config{
		Name:        appConfig.GetAgent().Name,
		Model:       llm,
		Instruction: instruction,
		Toolsets:    tools,
	})

}

func OfRunner(targetAgent agent.Agent, memoryService memory.Service, sessionService session.Service, withAutoCreateSession bool) (*runner.Runner, error) {
	return runner.New(runner.Config{
		AppName:           targetAgent.Name(),
		Agent:             targetAgent,
		SessionService:    sessionService,
		MemoryService:     memoryService,
		AutoCreateSession: withAutoCreateSession,
	})
}

func OfSkillBasedTool(ctx context.Context, location string) (*skilltoolset.SkillToolset, error) {
	return skilltoolset.New(ctx, skilltoolset.Config{
		Source: skill.NewFileSystemSource(os.DirFS(location)),
	})
}
