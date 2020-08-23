package encryption

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomBytes(t *testing.T) {
	b, _ := generateRandomBytes(1)
	require.NotNil(t, b)
}

func TestGenerateRandomBytesNotPanics(t *testing.T) {
	require.NotPanics(t, func() {
		generateRandomBytes(1) // nolint
	})
}

