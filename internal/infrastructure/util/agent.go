package util

import (
	"car-valuation-agent/internal/infrastructure/config"
	"context"
	"fmt"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/model/openaimodel"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/genai"
)

func OfAgent(ctx context.Context, instruction string, appConfig config.AppConfig, tools []tool.Toolset) (agent.Agent, error) {
	agentConfig := appConfig.GetAgent()
	modelConfig := appConfig.GetAgent().GetModel()
	modelAdvancedSetup := modelConfig.AdvanceSetup
	if modelAdvancedSetup.HasAdvancedSetup() {
		return withAdvanceSetup(ctx, instruction, agentConfig, modelConfig, modelAdvancedSetup, tools)
	} else {
		return withSimpleSetup(ctx, instruction, agentConfig, modelConfig, tools)
	}
}

func ofModel(ctx context.Context, modelConfig config.Model) (model.LLM, error) {
	llm, err := openaimodel.NewModel(ctx, modelConfig.ModelName, &openaimodel.ClientConfig{
		BaseURL: modelConfig.BaseUrl,
		APIKey:  modelConfig.ApiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("util: failed to create model: %w", err)
	}

	return llm, nil
}

func withSimpleSetup(ctx context.Context, instruction string, agentConfig config.Agent, modelConfig config.Model, tools []tool.Toolset) (agent.Agent, error) {
	llm, err := ofModel(ctx, modelConfig)
	if err != nil {
		return nil, err
	}

	createdAgent, err := llmagent.New(llmagent.Config{
		Name:        agentConfig.Name,
		Model:       llm,
		Instruction: instruction,
		Toolsets:    tools,
	})
	if err != nil {
		return nil, fmt.Errorf("util: failed to create agent: %w", err)
	}

	return createdAgent, nil
}

func withAdvanceSetup(ctx context.Context, instruction string, agentConfig config.Agent, modelConfig config.Model, advancedSetup config.AdvancedSetup, tools []tool.Toolset) (agent.Agent, error) {
	llm, err := ofModel(ctx, modelConfig)
	if err != nil {
		return nil, err
	}

	createdAgent, err := llmagent.New(llmagent.Config{
		Name:  agentConfig.Name,
		Model: llm,
		GenerateContentConfig: &genai.GenerateContentConfig{
			Temperature:      &advancedSetup.Temperature,
			TopP:             &advancedSetup.TopP,
			TopK:             &advancedSetup.TopK,
			PresencePenalty:  &advancedSetup.PresencePenalty,
			FrequencyPenalty: &advancedSetup.FrequencyPenalty,
			MaxOutputTokens:  advancedSetup.MaximumTokens,
		},
		Instruction: instruction,
		Toolsets:    tools,
	})
	if err != nil {
		return nil, fmt.Errorf("util: failed to create agent: %w", err)
	}

	return createdAgent, nil
}
