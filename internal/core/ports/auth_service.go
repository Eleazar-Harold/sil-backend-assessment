package ports

import (
	"context"
	"time"

	"silbackendassessment/internal/adapters/auth"
	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/oidc"
)

// AuthService defines the contract for authentication operations
type AuthService interface {
	// Traditional JWT-based authentication
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
	ValidateToken(tokenString string) (*auth.Claims, error)

	// OpenID Connect authentication
	GetOIDCAuthURL(ctx context.Context) (string, string, error) // Returns URL and state
	HandleOIDCCallback(ctx context.Context, code, state string) (*OIDCLoginResponse, error)
	ValidateOIDCToken(ctx context.Context, idToken string) (*oidc.OIDCUserInfo, error)
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// OIDCLoginResponse represents an OpenID Connect login response
type OIDCLoginResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	Customer     *domain.Customer `json:"customer"`
	ExpiresAt    time.Time        `json:"expires_at"`
	IsNewUser    bool             `json:"is_new_user"`
}
