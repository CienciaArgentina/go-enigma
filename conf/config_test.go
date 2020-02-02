package conf

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConfigShouldReturnNonEmptyConfig(t *testing.T) {
	config := New()

	require.Equal(t, "cienciaArgentinaDev", config.Database.User)
}
