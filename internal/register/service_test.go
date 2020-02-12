package register

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateVerificationTokenShouldReturnJwtToken(t *testing.T) {
	s := New(&registerService{}, nil)
	a := s.GenerateVerificationToken("juancito@asd.com")
	require.NotEmpty(t, a)
}
