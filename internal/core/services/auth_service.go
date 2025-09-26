package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"silbackendassessment/internal/adapters/auth"
	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/oidc"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo     ports.UserRepository
	customerRepo ports.CustomerRepository
	jwtManager   auth.JWTManagerInterface
	oidcProvider oidc.Provider
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo ports.UserRepository,
	customerRepo ports.CustomerRepository,
	jwtManager auth.JWTManagerInterface,
	oidcProvider oidc.Provider,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		customerRepo: customerRepo,
		jwtManager:   jwtManager,
		oidcProvider: oidcProvider,
	}
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *ports.LoginRequest) (*ports.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// In a real application, you would verify the password hash here
	// For now, we'll assume the password is correct

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &ports.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresAt:    time.Now().Add(6 * time.Hour), // 6 hours from now
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *ports.RegisterRequest) (*ports.LoginResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	user := &domain.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &ports.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresAt:    time.Now().Add(6 * time.Hour), // 6 hours from now
	}, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*ports.LoginResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, uuid.MustParse(claims.UserID))
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &ports.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
		ExpiresAt:    time.Now().Add(6 * time.Hour), // 6 hours from now
	}, nil
}

// ValidateToken validates an access token and returns user claims
func (s *AuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	return s.jwtManager.ValidateToken(tokenString)
}

// GetOIDCAuthURL generates the authorization URL for OIDC flow
func (s *AuthService) GetOIDCAuthURL(ctx context.Context) (string, string, error) {
	if s.oidcProvider == nil {
		return "", "", errors.New("OIDC provider not configured")
	}

	state, err := auth.GenerateState()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate state: %w", err)
	}

	authURL := s.oidcProvider.GetAuthURL(state)
	return authURL, state, nil
}

// HandleOIDCCallback handles the OIDC callback and creates/authenticates customer
func (s *AuthService) HandleOIDCCallback(ctx context.Context, code, state string) (*ports.OIDCLoginResponse, error) {
	if s.oidcProvider == nil {
		return nil, errors.New("OIDC provider not configured")
	}

	// Exchange code for tokens
	oidcToken, err := s.oidcProvider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from ID token
	userInfo, err := s.oidcProvider.ValidateIDToken(ctx, oidcToken.IDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate ID token: %w", err)
	}

	// Check if customer already exists
	customer, err := s.customerRepo.GetByEmail(ctx, userInfo.Email)
	if err != nil && err.Error() != "customer not found" {
		return nil, fmt.Errorf("failed to check existing customer: %w", err)
	}

	isNewUser := false
	if customer == nil {
		// Create new customer
		customer = &domain.Customer{
			ID:        uuid.New(),
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
			Email:     userInfo.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.customerRepo.Create(ctx, customer); err != nil {
			return nil, fmt.Errorf("failed to create customer: %w", err)
		}
		isNewUser = true
	}

	// Generate JWT tokens for the customer
	accessToken, err := s.jwtManager.GenerateToken(customer.ID.String(), customer.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(customer.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &ports.OIDCLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Customer:     customer,
		ExpiresAt:    time.Now().Add(6 * time.Hour),
		IsNewUser:    isNewUser,
	}, nil
}

// ValidateOIDCToken validates an OIDC ID token and returns user info
func (s *AuthService) ValidateOIDCToken(ctx context.Context, idToken string) (*oidc.OIDCUserInfo, error) {
	if s.oidcProvider == nil {
		return nil, errors.New("OIDC provider not configured")
	}

	return s.oidcProvider.ValidateIDToken(ctx, idToken)
}
