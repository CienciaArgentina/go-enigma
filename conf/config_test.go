package conf

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConfigShouldReturnNonEmptyConfig(t *testing.T) {
	config, err := New()

	require.NoError(t, err)
	require.NotEmpty(t, config.Database.Database)
	require.Equal(t, "prueba", config.Database.Password)
}
