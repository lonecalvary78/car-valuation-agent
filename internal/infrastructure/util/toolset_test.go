package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/agent"
)

type SimpleToolArgs struct {
	Parameter string `json:"parameter"`
}

type SimpleToolResult struct {
	Data string `json:"data"`
}

func TestCreatedSkillTool(t *testing.T) {
	skilltoolset, err := OfSkillBasedTool(t.Context(), "skills")
	require.NoError(t, err)
	require.NotEmpty(t, skilltoolset)
}

func TestFunctionalTool(t *testing.T) {
	simpleTool, err := OfFunctionalTool("simple-tool", func(ctx agent.Context, args SimpleToolArgs) (SimpleToolResult, error) {
		return SimpleToolResult{
			Data: args.Parameter,
		}, nil
	})

	require.NoError(t, err)
	require.NotEmpty(t, simpleTool)
	require.Equal(t, "simple-tool", simpleTool.Name())
}
