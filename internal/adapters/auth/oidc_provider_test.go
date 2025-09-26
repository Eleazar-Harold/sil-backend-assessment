package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateState(t *testing.T) {
	t.Run("Generate state successfully", func(t *testing.T) {
		state, err := GenerateState()

		assert.NoError(t, err)
		assert.NotEmpty(t, state)
		assert.Len(t, state, 44) // Base64 encoded 32 bytes = 44 characters
	})

	t.Run("Generate multiple states are different", func(t *testing.T) {
		state1, err1 := GenerateState()
		state2, err2 := GenerateState()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, state1, state2)
	})
}

func TestValidateState(t *testing.T) {
	t.Run("Validate matching states", func(t *testing.T) {
		state := "test-state-123"
		result := ValidateState(state, state)

		assert.True(t, result)
	})

	t.Run("Validate non-matching states", func(t *testing.T) {
		expectedState := "expected-state"
		actualState := "actual-state"
		result := ValidateState(expectedState, actualState)

		assert.False(t, result)
	})

	t.Run("Validate empty states", func(t *testing.T) {
		result := ValidateState("", "")

		assert.True(t, result)
	})

	t.Run("Validate one empty state", func(t *testing.T) {
		result := ValidateState("", "non-empty")

		assert.False(t, result)
	})
}

func TestNewOIDCProvider_ErrorCases(t *testing.T) {
	t.Run("Create OIDC provider with invalid URL", func(t *testing.T) {
		provider, err := NewOIDCProvider(
			"invalid-url",
			"test-client-id",
			"test-client-secret",
			"http://localhost:8080/callback",
			[]string{"openid"},
		)

		assert.Error(t, err)
		assert.Nil(t, provider)
		assert.Contains(t, err.Error(), "failed to create OIDC provider")
	})

	t.Run("Create OIDC provider with empty URL", func(t *testing.T) {
		provider, err := NewOIDCProvider(
			"",
			"test-client-id",
			"test-client-secret",
			"http://localhost:8080/callback",
			[]string{"openid"},
		)

		assert.Error(t, err)
		assert.Nil(t, provider)
		assert.Contains(t, err.Error(), "failed to create OIDC provider")
	})
}
