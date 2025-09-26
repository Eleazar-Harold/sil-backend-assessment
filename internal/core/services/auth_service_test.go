package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"silbackendassessment/internal/adapters/auth"
	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/oidc"
	"silbackendassessment/internal/core/ports"
	"silbackendassessment/internal/testutils"

	"github.com/google/uuid"
)

// MockJWTManager for testing
type MockJWTManager struct {
	GenerateTokenFunc        func(userID, email string) (string, error)
	GenerateRefreshTokenFunc func(userID string) (string, error)
	ValidateTokenFunc        func(token string) (*auth.Claims, error)
	ValidateRefreshTokenFunc func(token string) (*auth.Claims, error)
}

// Ensure MockJWTManager implements auth.JWTManagerInterface
var _ auth.JWTManagerInterface = (*MockJWTManager)(nil)

func (m *MockJWTManager) GenerateToken(userID, email string) (string, error) {
	if m.GenerateTokenFunc != nil {
		return m.GenerateTokenFunc(userID, email)
	}
	return "mock-token", nil
}

func (m *MockJWTManager) GenerateRefreshToken(userID string) (string, error) {
	if m.GenerateRefreshTokenFunc != nil {
		return m.GenerateRefreshTokenFunc(userID)
	}
	return "mock-refresh-token", nil
}

func (m *MockJWTManager) ValidateToken(token string) (*auth.Claims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(token)
	}
	return &auth.Claims{
		UserID: "mock-user-id",
		Email:  "test@example.com",
	}, nil
}

func (m *MockJWTManager) ValidateRefreshToken(token string) (*auth.Claims, error) {
	if m.ValidateRefreshTokenFunc != nil {
		return m.ValidateRefreshTokenFunc(token)
	}
	return &auth.Claims{
		UserID: "mock-user-id",
		Email:  "test@example.com",
	}, nil
}

// MockOIDCProvider for testing
type MockOIDCProvider struct {
	GetAuthURLFunc      func(state string) string
	ExchangeCodeFunc    func(ctx context.Context, code string) (*oidc.OIDCToken, error)
	GetUserInfoFunc     func(ctx context.Context, accessToken string) (*oidc.OIDCUserInfo, error)
	ValidateIDTokenFunc func(ctx context.Context, idToken string) (*oidc.OIDCUserInfo, error)
	RefreshTokenFunc    func(ctx context.Context, refreshToken string) (*oidc.OIDCToken, error)
}

func (m *MockOIDCProvider) GetAuthURL(state string) string {
	if m.GetAuthURLFunc != nil {
		return m.GetAuthURLFunc(state)
	}
	return "https://mock-provider.com/auth"
}

func (m *MockOIDCProvider) ExchangeCode(ctx context.Context, code string) (*oidc.OIDCToken, error) {
	if m.ExchangeCodeFunc != nil {
		return m.ExchangeCodeFunc(ctx, code)
	}
	return &oidc.OIDCToken{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}, nil
}

func (m *MockOIDCProvider) GetUserInfo(ctx context.Context, accessToken string) (*oidc.OIDCUserInfo, error) {
	if m.GetUserInfoFunc != nil {
		return m.GetUserInfoFunc(ctx, accessToken)
	}
	return &oidc.OIDCUserInfo{
		Subject: "mock-sub",
		Email:   "test@example.com",
		Name:    "Test User",
	}, nil
}

func (m *MockOIDCProvider) ValidateIDToken(ctx context.Context, idToken string) (*oidc.OIDCUserInfo, error) {
	if m.ValidateIDTokenFunc != nil {
		return m.ValidateIDTokenFunc(ctx, idToken)
	}
	return &oidc.OIDCUserInfo{
		Subject: "mock-sub",
		Email:   "test@example.com",
		Name:    "Test User",
	}, nil
}

func (m *MockOIDCProvider) RefreshToken(ctx context.Context, refreshToken string) (*oidc.OIDCToken, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, refreshToken)
	}
	return &oidc.OIDCToken{
		AccessToken: "mock-access-token",
		ExpiresAt:   time.Now().Add(time.Hour),
	}, nil
}

func TestAuthService_Login(t *testing.T) {
	mockUserRepo := testutils.NewMockUserRepository()
	mockCustomerRepo := testutils.NewMockCustomerRepository()
	mockJWTManager := &MockJWTManager{}
	mockOIDCProvider := &MockOIDCProvider{}
	service := NewAuthService(mockUserRepo, mockCustomerRepo, mockJWTManager, mockOIDCProvider)
	ctx := context.Background()

	t.Run("Login successfully", func(t *testing.T) {
		// Set up user
		userID := uuid.New()
		user := &domain.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockUserRepo.UsersByEmail["john@example.com"] = user

		// Set up JWT manager
		mockJWTManager.GenerateTokenFunc = func(id, email string) (string, error) {
			if id != userID.String() || email != "john@example.com" {
				t.Errorf("Expected userID %s and email john@example.com, got: %s, %s", userID, id, email)
			}
			return "jwt-token", nil
		}

		req := &ports.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		response, err := service.Login(ctx, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if response == nil {
			t.Error("Expected response to be returned")
		}

		if response.AccessToken != "jwt-token" {
			t.Errorf("Expected AccessToken to be 'jwt-token', got: %s", response.AccessToken)
		}

		if response.User.ID != userID {
			t.Errorf("Expected User ID to be %s, got: %s", userID, response.User.ID)
		}
	})

	t.Run("Login with non-existent user", func(t *testing.T) {
		req := &ports.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		response, err := service.Login(ctx, req)

		if err == nil {
			t.Error("Expected error for non-existent user")
		}

		if response != nil {
			t.Error("Expected response to be nil")
		}

		if err.Error() != "invalid credentials" {
			t.Errorf("Expected 'invalid credentials' error, got: %v", err)
		}
	})

	t.Run("Login with repository error", func(t *testing.T) {
		mockUserRepo.GetByEmailError = errors.New("database error")

		req := &ports.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		response, err := service.Login(ctx, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if response != nil {
			t.Error("Expected response to be nil")
		}

		if err.Error() != "invalid credentials" {
			t.Errorf("Expected 'invalid credentials' error, got: %v", err)
		}
	})

	t.Run("Login with JWT generation error", func(t *testing.T) {
		// Create fresh mocks for this test
		mockUserRepo := testutils.NewMockUserRepository()
		mockCustomerRepo := testutils.NewMockCustomerRepository()
		mockJWTManager := &MockJWTManager{}
		mockOIDCProvider := &MockOIDCProvider{}
		service := NewAuthService(mockUserRepo, mockCustomerRepo, mockJWTManager, mockOIDCProvider)

		// Set up user
		userID := uuid.New()
		user := &domain.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockUserRepo.UsersByEmail["john@example.com"] = user

		// Set up JWT manager to return error
		mockJWTManager.GenerateTokenFunc = func(userID string, email string) (string, error) {
			return "", errors.New("JWT generation failed")
		}

		req := &ports.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		response, err := service.Login(ctx, req)

		if err == nil {
			t.Error("Expected error from JWT generation")
		}

		if response != nil {
			t.Error("Expected response to be nil")
		}

		if err.Error() != "failed to generate access token" {
			t.Errorf("Expected JWT generation error, got: %v", err)
		}
	})
}

func TestAuthService_Register(t *testing.T) {
	mockUserRepo := testutils.NewMockUserRepository()
	mockCustomerRepo := testutils.NewMockCustomerRepository()
	mockJWTManager := &MockJWTManager{}
	mockOIDCProvider := &MockOIDCProvider{}
	service := NewAuthService(mockUserRepo, mockCustomerRepo, mockJWTManager, mockOIDCProvider)
	ctx := context.Background()

	t.Run("Register successfully", func(t *testing.T) {
		req := &ports.RegisterRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		// Set up JWT manager
		mockJWTManager.GenerateTokenFunc = func(id, email string) (string, error) {
			return "jwt-token", nil
		}

		response, err := service.Register(ctx, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if response == nil {
			t.Error("Expected response to be returned")
		}

		if response.AccessToken != "jwt-token" {
			t.Errorf("Expected AccessToken to be 'jwt-token', got: %s", response.AccessToken)
		}

		if response.User.Name != "John Doe" {
			t.Errorf("Expected User Name to be 'John Doe', got: %s", response.User.Name)
		}

		if response.User.Email != "john@example.com" {
			t.Errorf("Expected User Email to be 'john@example.com', got: %s", response.User.Email)
		}
	})

	t.Run("Register with existing email", func(t *testing.T) {
		// Set up existing user
		existingUser := &domain.User{
			ID:        uuid.New(),
			Name:      "Existing User",
			Email:     "existing@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockUserRepo.UsersByEmail["existing@example.com"] = existingUser

		req := &ports.RegisterRequest{
			Name:     "New User",
			Email:    "existing@example.com",
			Password: "password123",
		}

		response, err := service.Register(ctx, req)

		if err == nil {
			t.Error("Expected error for existing email")
		}

		if response != nil {
			t.Error("Expected response to be nil")
		}

		if err.Error() != "user with this email already exists" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("Register with repository error", func(t *testing.T) {
		// Create fresh mocks for this test
		mockUserRepo := testutils.NewMockUserRepository()
		mockCustomerRepo := testutils.NewMockCustomerRepository()
		mockJWTManager := &MockJWTManager{}
		mockOIDCProvider := &MockOIDCProvider{}
		service := NewAuthService(mockUserRepo, mockCustomerRepo, mockJWTManager, mockOIDCProvider)

		mockUserRepo.CreateError = errors.New("database error")

		req := &ports.RegisterRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		response, err := service.Register(ctx, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if response != nil {
			t.Error("Expected response to be nil")
		}

		if err.Error() != "failed to create user" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	mockUserRepo := testutils.NewMockUserRepository()
	mockCustomerRepo := testutils.NewMockCustomerRepository()
	mockJWTManager := &MockJWTManager{}
	mockOIDCProvider := &MockOIDCProvider{}
	service := NewAuthService(mockUserRepo, mockCustomerRepo, mockJWTManager, mockOIDCProvider)

	t.Run("Validate token successfully", func(t *testing.T) {
		userID := uuid.New()
		user := &domain.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockUserRepo.Users[userID] = user

		// Set up JWT manager
		mockJWTManager.ValidateTokenFunc = func(token string) (*auth.Claims, error) {
			if token != "valid-token" {
				t.Errorf("Expected token 'valid-token', got: %s", token)
			}
			return &auth.Claims{
				UserID: userID.String(),
				Email:  "john@example.com",
			}, nil
		}

		claims, err := service.ValidateToken("valid-token")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if claims == nil {
			t.Error("Expected claims to be returned")
		}

		if claims.UserID != userID.String() {
			t.Errorf("Expected User ID to be %s, got: %s", userID, claims.UserID)
		}
	})

	t.Run("Validate invalid token", func(t *testing.T) {
		// Set up JWT manager to return error
		mockJWTManager.ValidateTokenFunc = func(token string) (*auth.Claims, error) {
			return nil, errors.New("invalid token")
		}

		claims, err := service.ValidateToken("invalid-token")

		if err == nil {
			t.Error("Expected error for invalid token")
		}

		if claims != nil {
			t.Error("Expected claims to be nil")
		}

		if err.Error() != "invalid token" {
			t.Errorf("Expected 'invalid token' error, got: %v", err)
		}
	})

}
