package oidc

import "context"

// Provider defines the contract for OpenID Connect operations
type Provider interface {
	// GetAuthURL generates the authorization URL for OIDC flow
	GetAuthURL(state string) string

	// ExchangeCode exchanges authorization code for tokens
	ExchangeCode(ctx context.Context, code string) (*OIDCToken, error)

	// GetUserInfo retrieves user information using access token
	GetUserInfo(ctx context.Context, accessToken string) (*OIDCUserInfo, error)

	// ValidateIDToken validates and parses the ID token
	ValidateIDToken(ctx context.Context, idToken string) (*OIDCUserInfo, error)

	// RefreshToken refreshes the access token using refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*OIDCToken, error)
}

