package util

import (
	"car-valuation-agent/internal/infrastructure/config"
	"context"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model/openaimodel"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/genai"
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

func OfAgentWithAdvancedConfiguration(
	ctx context.Context,
	instruction string,
	appConfig config.AppConfig,
	tools []tool.Toolset,
	temperature float32,
	topP float32,
	topK float32,
	presencePenalty float32,
	frequencyPenalty float32,
	maximumOutputTokens int32) (agent.Agent, error) {
	llm, err := openaimodel.NewModel(ctx, appConfig.GetAgent().GetModel().ModelName, &openaimodel.ClientConfig{
		BaseURL: appConfig.GetAgent().GetModel().BaseUrl,
		APIKey:  appConfig.GetAgent().GetModel().ApiKey,
	})
	if err != nil {
		return nil, err
	}

	return llmagent.New(llmagent.Config{
		Name:  appConfig.GetAgent().Name,
		Model: llm,
		GenerateContentConfig: &genai.GenerateContentConfig{
			Temperature:      &temperature,
			TopP:             &topP,
			TopK:             &topK,
			PresencePenalty:  &presencePenalty,
			FrequencyPenalty: &frequencyPenalty,
			MaxOutputTokens:  maximumOutputTokens,
		},
		Instruction: instruction,
		Toolsets:    tools,
	})
}
