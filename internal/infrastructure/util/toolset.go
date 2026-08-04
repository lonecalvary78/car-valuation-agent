package util

import (
	"car-valuation-agent/internal/infrastructure/validator"
	"context"
	"os"

	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/adk/v2/tool/skilltoolset"
	"google.golang.org/adk/v2/tool/skilltoolset/skill"
)

func OfSkillBasedTool(ctx context.Context, location string) (*skilltoolset.SkillToolset, error) {
	if err := validator.ValidateForRequiredOfString("Skill Location", location); err != nil {
		return nil, err
	}
	return skilltoolset.New(ctx, skilltoolset.Config{
		Source: skill.NewFileSystemSource(os.DirFS(location)),
	})
}

func OfFunctionalTool[TArgs, TResults any](toolName string, calledFunction functiontool.Func[TArgs, TResults]) (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name: toolName,
	}, calledFunction)
}
