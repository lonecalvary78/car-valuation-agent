package util

import (
	"car-valuation-agent/internal/infrastructure/validator"
	"context"
	"fmt"
	"os"

	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/adk/v2/tool/skilltoolset"
	"google.golang.org/adk/v2/tool/skilltoolset/skill"
)

func OfSkillBasedTool(ctx context.Context, location string) (*skilltoolset.SkillToolset, error) {
	err := validator.ValidateForRequiredOfString("Skill Location", location)
	if err != nil {
		return nil, fmt.Errorf("util: invalid skill location: %w", err)
	}

	skillToolset, err := skilltoolset.New(ctx, skilltoolset.Config{
		Source: skill.NewFileSystemSource(os.DirFS(location)),
	})
	if err != nil {
		return nil, fmt.Errorf("util: failed to create skill toolset: %w", err)
	}

	return skillToolset, nil
}

func OfFunctionalTool[TArgs, TResults any](toolName string, calledFunction functiontool.Func[TArgs, TResults]) (tool.Tool, error) {
	createdTool, err := functiontool.New(functiontool.Config{
		Name: toolName,
	}, calledFunction)
	if err != nil {
		return nil, fmt.Errorf("util: failed to create functional tool: %w", err)
	}

	return createdTool, nil
}
