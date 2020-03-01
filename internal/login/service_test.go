package login

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultLoginOptionsShouldReturnDefaultOptions(t *testing.T) {
	opt := defaultLoginOptions()
	require.NotNil(t, opt.LockoutOptions.LockoutTimeDuration)
}
