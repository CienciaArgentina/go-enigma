package encryption

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateRandomBytes(t *testing.T) {
	b, _ := GenerateRandomBytes(1)
	require.NotNil(t, b)
}
