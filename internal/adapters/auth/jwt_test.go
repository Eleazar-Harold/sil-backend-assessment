package auth

import (
	"testing"
	"time"
)

func TestNewJWTManager(t *testing.T) {
	secretKey := "test-secret-key"
	refreshSecret := "test-refresh-secret"
	tokenExpiry := 1 * time.Hour
	refreshExpiry := 24 * time.Hour

	manager := NewJWTManager(secretKey, refreshSecret, tokenExpiry, refreshExpiry)

	if manager == nil {
		t.Error("Expected JWTManager to be created")
	}

	if manager.secretKey != secretKey {
		t.Errorf("Expected secretKey to be %s, got: %s", secretKey, manager.secretKey)
	}

	if manager.refreshSecret != refreshSecret {
		t.Errorf("Expected refreshSecret to be %s, got: %s", refreshSecret, manager.refreshSecret)
	}

	if manager.tokenExpiry != tokenExpiry {
		t.Errorf("Expected tokenExpiry to be %v, got: %v", tokenExpiry, manager.tokenExpiry)
	}

	if manager.refreshExpiry != refreshExpiry {
		t.Errorf("Expected refreshExpiry to be %v, got: %v", refreshExpiry, manager.refreshExpiry)
	}
}

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-refresh-secret", time.Hour, 24*time.Hour)

	t.Run("Generate valid token", func(t *testing.T) {
		userID := "user-123"
		email := "test@example.com"

		token, err := manager.GenerateToken(userID, email)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if token == "" {
			t.Error("Expected token to be generated")
		}

		// Token should be a JWT with 3 parts separated by dots
		parts := len([]rune(token))
		if parts < 100 { // JWT tokens are much longer than 100 characters
			t.Error("Expected token to be a valid JWT format")
		}
	})

	t.Run("Generate token with empty userID", func(t *testing.T) {
		token, err := manager.GenerateToken("", "test@example.com")

		if err != nil {
			t.Errorf("Expected no error for empty userID, got: %v", err)
		}

		if token == "" {
			t.Error("Expected token to be generated even with empty userID")
		}
	})

	t.Run("Generate token with empty email", func(t *testing.T) {
		token, err := manager.GenerateToken("user-123", "")

		if err != nil {
			t.Errorf("Expected no error for empty email, got: %v", err)
		}

		if token == "" {
			t.Error("Expected token to be generated even with empty email")
		}
	})
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-refresh-secret", time.Hour, 24*time.Hour)

	t.Run("Generate valid refresh token", func(t *testing.T) {
		userID := "user-123"

		token, err := manager.GenerateRefreshToken(userID)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if token == "" {
			t.Error("Expected refresh token to be generated")
		}

		// Token should be a JWT with 3 parts separated by dots
		parts := len([]rune(token))
		if parts < 100 { // JWT tokens are much longer than 100 characters
			t.Error("Expected refresh token to be a valid JWT format")
		}
	})

	t.Run("Generate refresh token with empty userID", func(t *testing.T) {
		token, err := manager.GenerateRefreshToken("")

		if err != nil {
			t.Errorf("Expected no error for empty userID, got: %v", err)
		}

		if token == "" {
			t.Error("Expected refresh token to be generated even with empty userID")
		}
	})
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-refresh-secret", time.Hour, 24*time.Hour)

	t.Run("Validate valid token", func(t *testing.T) {
		userID := "user-123"
		email := "test@example.com"

		token, err := manager.GenerateToken(userID, email)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		claims, err := manager.ValidateToken(token)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if claims == nil {
			t.Error("Expected claims to be returned")
		}

		if claims.UserID != userID {
			t.Errorf("Expected UserID to be %s, got: %s", userID, claims.UserID)
		}

		if claims.Email != email {
			t.Errorf("Expected Email to be %s, got: %s", email, claims.Email)
		}
	})

	t.Run("Validate invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"

		claims, err := manager.ValidateToken(invalidToken)

		if err == nil {
			t.Error("Expected error for invalid token")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for invalid token")
		}
	})

	t.Run("Validate empty token", func(t *testing.T) {
		claims, err := manager.ValidateToken("")

		if err == nil {
			t.Error("Expected error for empty token")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for empty token")
		}
	})

	t.Run("Validate token with wrong secret", func(t *testing.T) {
		// Create a token with one manager
		manager1 := NewJWTManager("secret1", "refresh1", time.Hour, 24*time.Hour)
		userID := "user-123"
		email := "test@example.com"

		token, err := manager1.GenerateToken(userID, email)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Try to validate with a different manager (different secret)
		manager2 := NewJWTManager("secret2", "refresh2", time.Hour, 24*time.Hour)
		claims, err := manager2.ValidateToken(token)

		if err == nil {
			t.Error("Expected error for token with wrong secret")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for token with wrong secret")
		}
	})
}

func TestJWTManager_ValidateRefreshToken(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-refresh-secret", time.Hour, 24*time.Hour)

	t.Run("Validate valid refresh token", func(t *testing.T) {
		userID := "user-123"

		token, err := manager.GenerateRefreshToken(userID)
		if err != nil {
			t.Fatalf("Failed to generate refresh token: %v", err)
		}

		claims, err := manager.ValidateRefreshToken(token)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if claims == nil {
			t.Error("Expected claims to be returned")
		}

		if claims.UserID != userID {
			t.Errorf("Expected UserID to be %s, got: %s", userID, claims.UserID)
		}
	})

	t.Run("Validate invalid refresh token", func(t *testing.T) {
		invalidToken := "invalid.refresh.token"

		claims, err := manager.ValidateRefreshToken(invalidToken)

		if err == nil {
			t.Error("Expected error for invalid refresh token")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for invalid refresh token")
		}
	})

	t.Run("Validate empty refresh token", func(t *testing.T) {
		claims, err := manager.ValidateRefreshToken("")

		if err == nil {
			t.Error("Expected error for empty refresh token")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for empty refresh token")
		}
	})

	t.Run("Validate refresh token with wrong secret", func(t *testing.T) {
		// Create a refresh token with one manager
		manager1 := NewJWTManager("secret1", "refresh1", time.Hour, 24*time.Hour)
		userID := "user-123"

		token, err := manager1.GenerateRefreshToken(userID)
		if err != nil {
			t.Fatalf("Failed to generate refresh token: %v", err)
		}

		// Try to validate with a different manager (different refresh secret)
		manager2 := NewJWTManager("secret2", "refresh2", time.Hour, 24*time.Hour)
		claims, err := manager2.ValidateRefreshToken(token)

		if err == nil {
			t.Error("Expected error for refresh token with wrong secret")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for refresh token with wrong secret")
		}
	})

	t.Run("Validate access token as refresh token", func(t *testing.T) {
		userID := "user-123"
		email := "test@example.com"

		// Generate access token
		accessToken, err := manager.GenerateToken(userID, email)
		if err != nil {
			t.Fatalf("Failed to generate access token: %v", err)
		}

		// Try to validate access token as refresh token
		claims, err := manager.ValidateRefreshToken(accessToken)

		if err == nil {
			t.Error("Expected error when validating access token as refresh token")
		}

		if claims != nil {
			t.Error("Expected claims to be nil when validating access token as refresh token")
		}
	})
}

func TestJWTManager_ValidateToken_EdgeCases(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-refresh-secret", time.Hour, 24*time.Hour)

	t.Run("Validate token with invalid signature", func(t *testing.T) {
		// Create a token with different secret
		differentManager := NewJWTManager("different-secret", "test-refresh-secret", time.Hour, 24*time.Hour)
		userID := "user-123"
		email := "test@example.com"

		token, err := differentManager.GenerateToken(userID, email)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Try to validate with original manager
		claims, err := manager.ValidateToken(token)

		if err == nil {
			t.Error("Expected error for token with invalid signature")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for invalid token")
		}
	})

	t.Run("Validate malformed token", func(t *testing.T) {
		malformedToken := "not.a.valid.jwt.token"

		claims, err := manager.ValidateToken(malformedToken)

		if err == nil {
			t.Error("Expected error for malformed token")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for malformed token")
		}
	})

	t.Run("Validate empty token", func(t *testing.T) {
		claims, err := manager.ValidateToken("")

		if err == nil {
			t.Error("Expected error for empty token")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for empty token")
		}
	})

	t.Run("Validate token with missing parts", func(t *testing.T) {
		invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" // Only header

		claims, err := manager.ValidateToken(invalidToken)

		if err == nil {
			t.Error("Expected error for token with missing parts")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for invalid token")
		}
	})
}

func TestJWTManager_ValidateRefreshToken_EdgeCases(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-refresh-secret", time.Hour, 24*time.Hour)

	t.Run("Validate refresh token with invalid signature", func(t *testing.T) {
		// Create a refresh token with different secret
		differentManager := NewJWTManager("test-secret", "different-refresh-secret", time.Hour, 24*time.Hour)
		userID := "user-123"

		refreshToken, err := differentManager.GenerateRefreshToken(userID)
		if err != nil {
			t.Fatalf("Failed to generate refresh token: %v", err)
		}

		// Try to validate with original manager
		claims, err := manager.ValidateRefreshToken(refreshToken)

		if err == nil {
			t.Error("Expected error for refresh token with invalid signature")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for invalid refresh token")
		}
	})

	t.Run("Validate malformed refresh token", func(t *testing.T) {
		malformedToken := "not.a.valid.refresh.token"

		claims, err := manager.ValidateRefreshToken(malformedToken)

		if err == nil {
			t.Error("Expected error for malformed refresh token")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for malformed refresh token")
		}
	})

	t.Run("Validate empty refresh token", func(t *testing.T) {
		claims, err := manager.ValidateRefreshToken("")

		if err == nil {
			t.Error("Expected error for empty refresh token")
		}

		if claims != nil {
			t.Error("Expected claims to be nil for empty refresh token")
		}
	})
}

func TestJWTManager_GenerateToken_EdgeCases(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-refresh-secret", time.Hour, 24*time.Hour)

	t.Run("Generate token with empty userID", func(t *testing.T) {
		token, err := manager.GenerateToken("", "test@example.com")

		if err != nil {
			t.Errorf("Expected no error for empty userID, got: %v", err)
		}

		if token == "" {
			t.Error("Expected token to be generated even with empty userID")
		}
	})

	t.Run("Generate token with empty email", func(t *testing.T) {
		token, err := manager.GenerateToken("user-123", "")

		if err != nil {
			t.Errorf("Expected no error for empty email, got: %v", err)
		}

		if token == "" {
			t.Error("Expected token to be generated even with empty email")
		}
	})

	t.Run("Generate token with special characters in userID", func(t *testing.T) {
		userID := "user-123!@#$%^&*()"
		email := "test@example.com"

		token, err := manager.GenerateToken(userID, email)

		if err != nil {
			t.Errorf("Expected no error for special characters in userID, got: %v", err)
		}

		if token == "" {
			t.Error("Expected token to be generated with special characters")
		}

		// Validate the token
		claims, err := manager.ValidateToken(token)
		if err != nil {
			t.Errorf("Expected token to be valid, got error: %v", err)
		}

		if claims.UserID != userID {
			t.Errorf("Expected UserID %s, got: %s", userID, claims.UserID)
		}
	})

	t.Run("Generate token with special characters in email", func(t *testing.T) {
		userID := "user-123"
		email := "test+tag@example.com"

		token, err := manager.GenerateToken(userID, email)

		if err != nil {
			t.Errorf("Expected no error for special characters in email, got: %v", err)
		}

		if token == "" {
			t.Error("Expected token to be generated with special characters in email")
		}

		// Validate the token
		claims, err := manager.ValidateToken(token)
		if err != nil {
			t.Errorf("Expected token to be valid, got error: %v", err)
		}

		if claims.Email != email {
			t.Errorf("Expected Email %s, got: %s", email, claims.Email)
		}
	})
}

func TestJWTManager_GenerateRefreshToken_EdgeCases(t *testing.T) {
	manager := NewJWTManager("test-secret", "test-refresh-secret", time.Hour, 24*time.Hour)

	t.Run("Generate refresh token with empty userID", func(t *testing.T) {
		refreshToken, err := manager.GenerateRefreshToken("")

		if err != nil {
			t.Errorf("Expected no error for empty userID, got: %v", err)
		}

		if refreshToken == "" {
			t.Error("Expected refresh token to be generated even with empty userID")
		}
	})

	t.Run("Generate refresh token with special characters", func(t *testing.T) {
		userID := "user-123!@#$%^&*()"

		refreshToken, err := manager.GenerateRefreshToken(userID)

		if err != nil {
			t.Errorf("Expected no error for special characters in userID, got: %v", err)
		}

		if refreshToken == "" {
			t.Error("Expected refresh token to be generated with special characters")
		}

		// Validate the refresh token
		claims, err := manager.ValidateRefreshToken(refreshToken)
		if err != nil {
			t.Errorf("Expected refresh token to be valid, got error: %v", err)
		}

		if claims.UserID != userID {
			t.Errorf("Expected UserID %s, got: %s", userID, claims.UserID)
		}
	})
}
