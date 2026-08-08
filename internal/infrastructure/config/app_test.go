package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	SetEnv(t)
	appConfig := Load()
	err := appConfig.Validate()
	require.NoError(t, err)
}
