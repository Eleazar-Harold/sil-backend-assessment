package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"silbackendassessment/internal/core/oidc"

	oidclib "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// oidcProvider implements the OIDCProvider interface
type oidcProvider struct {
	provider     *oidclib.Provider
	verifier     *oidclib.IDTokenVerifier
	oauth2Config *oauth2.Config
}

// NewOIDCProvider creates a new OpenID Connect provider
func NewOIDCProvider(providerURL, clientID, clientSecret, redirectURL string, scopes []string) (oidc.Provider, error) {
	ctx := context.Background()

	provider, err := oidclib.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidclib.Config{
		ClientID: clientID,
	})

	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	return &oidcProvider{
		provider:     provider,
		verifier:     verifier,
		oauth2Config: oauth2Config,
	}, nil
}

// GetAuthURL generates the authorization URL for OIDC flow
func (p *oidcProvider) GetAuthURL(state string) string {
	return p.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges authorization code for tokens
func (p *oidcProvider) ExchangeCode(ctx context.Context, code string) (*oidc.OIDCToken, error) {
	token, err := p.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Extract ID token from the response
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token in response")
	}

	return &oidc.OIDCToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      idToken,
		ExpiresAt:    token.Expiry,
		TokenType:    token.TokenType,
	}, nil
}

// GetUserInfo retrieves user information using access token
func (p *oidcProvider) GetUserInfo(ctx context.Context, accessToken string) (*oidc.OIDCUserInfo, error) {
	// Get userinfo URL from provider
	userInfoURL := p.provider.Endpoint().AuthURL
	// Replace /authorize with /userinfo for most providers
	if strings.Contains(userInfoURL, "/authorize") {
		userInfoURL = strings.Replace(userInfoURL, "/authorize", "/userinfo", 1)
	} else {
		// Fallback to well-known userinfo endpoint
		userInfoURL = strings.TrimSuffix(p.provider.Endpoint().AuthURL, "/authorize") + "/userinfo"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed with status: %d", resp.StatusCode)
	}

	var userInfo oidc.OIDCUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo response: %w", err)
	}

	return &userInfo, nil
}

// ValidateIDToken validates and parses the ID token
func (p *oidcProvider) ValidateIDToken(ctx context.Context, idToken string) (*oidc.OIDCUserInfo, error) {
	token, err := p.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	var claims struct {
		Subject   string `json:"sub"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		FirstName string `json:"given_name"`
		LastName  string `json:"family_name"`
		Picture   string `json:"picture"`
	}

	if err := token.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse ID token claims: %w", err)
	}

	return &oidc.OIDCUserInfo{
		Subject:   claims.Subject,
		Email:     claims.Email,
		Name:      claims.Name,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
		Picture:   claims.Picture,
	}, nil
}

// RefreshToken refreshes the access token using refresh token
func (p *oidcProvider) RefreshToken(ctx context.Context, refreshToken string) (*oidc.OIDCToken, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := p.oauth2Config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Extract ID token from the response
	idToken, ok := newToken.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token in refresh response")
	}

	return &oidc.OIDCToken{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		IDToken:      idToken,
		ExpiresAt:    newToken.Expiry,
		TokenType:    newToken.TokenType,
	}, nil
}

// GenerateState generates a random state parameter for OAuth2 flow
func GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateState validates the state parameter
func ValidateState(expected, actual string) bool {
	return expected == actual
}
