package oidc

import "time"

// OIDCUserInfo represents user information from OpenID Connect provider
type OIDCUserInfo struct {
	Subject   string `json:"sub"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Picture   string `json:"picture"`
}

// OIDCToken represents tokens from OpenID Connect provider
type OIDCToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	IDToken      string    `json:"id_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

