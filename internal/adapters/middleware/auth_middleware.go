package middleware

import (
	"context"
	"net/http"

	"silbackendassessment/internal/core/ports"

	"github.com/uptrace/bunrouter"
)

// AuthMiddleware handles authentication for both JWT and OIDC tokens
type AuthMiddleware struct {
	authService ports.AuthService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService ports.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// CustomerContextKey is the key used to store customer info in context
type CustomerContextKey struct{}

// UserContextKey is the key used to store user info in context
type UserContextKey struct{}

// CustomerInfo represents customer information stored in context
type CustomerInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// UserInfo represents user information stored in context
type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// RequireAuth middleware that validates JWT tokens
func (m *AuthMiddleware) RequireAuth(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		// Get token from Authorization header
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return nil
		}

		// Extract token from "Bearer <token>" format
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return nil
		}

		token := authHeader[7:]

		// Validate the JWT token
		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return err
		}

		// Add user info to context
		userInfo := &UserInfo{
			ID:    claims.UserID,
			Email: claims.Email,
		}

		ctx := context.WithValue(req.Context(), UserContextKey{}, userInfo)
		req = req.WithContext(ctx)

		return next(w, req)
	}
}

// RequireCustomerAuth middleware that validates tokens and ensures customer context
func (m *AuthMiddleware) RequireCustomerAuth(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		// Get token from Authorization header
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return nil
		}

		// Extract token from "Bearer <token>" format
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return nil
		}

		token := authHeader[7:]

		// Try to validate as JWT token first
		claims, err := m.authService.ValidateToken(token)
		if err == nil {
			// JWT token is valid, add customer info to context
			customerInfo := &CustomerInfo{
				ID:    claims.UserID,
				Email: claims.Email,
			}

			ctx := context.WithValue(req.Context(), CustomerContextKey{}, customerInfo)
			req = req.WithContext(ctx)

			return next(w, req)
		}

		// If JWT validation fails, try OIDC token validation
		userInfo, err := m.authService.ValidateOIDCToken(req.Context(), token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return err
		}

		// OIDC token is valid, add customer info to context
		customerInfo := &CustomerInfo{
			ID:    userInfo.Subject,
			Email: userInfo.Email,
		}

		ctx := context.WithValue(req.Context(), CustomerContextKey{}, customerInfo)
		req = req.WithContext(ctx)

		return next(w, req)
	}
}

// RequireOIDCAuth middleware that specifically validates OIDC tokens
func (m *AuthMiddleware) RequireOIDCAuth(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		// Get token from Authorization header
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return nil
		}

		// Extract token from "Bearer <token>" format
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return nil
		}

		token := authHeader[7:]

		// Validate the OIDC token
		userInfo, err := m.authService.ValidateOIDCToken(req.Context(), token)
		if err != nil {
			http.Error(w, "Invalid OIDC token: "+err.Error(), http.StatusUnauthorized)
			return err
		}

		// Add customer info to context
		customerInfo := &CustomerInfo{
			ID:    userInfo.Subject,
			Email: userInfo.Email,
		}

		ctx := context.WithValue(req.Context(), CustomerContextKey{}, customerInfo)
		req = req.WithContext(ctx)

		return next(w, req)
	}
}

// OptionalAuth middleware that validates tokens if present but doesn't require them
func (m *AuthMiddleware) OptionalAuth(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		// Get token from Authorization header
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			return next(w, req)
		}

		// Extract token from "Bearer <token>" format
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			// Invalid format, continue without authentication
			return next(w, req)
		}

		token := authHeader[7:]

		// Try to validate as JWT token first
		claims, err := m.authService.ValidateToken(token)
		if err == nil {
			// JWT token is valid, add user info to context
			userInfo := &UserInfo{
				ID:    claims.UserID,
				Email: claims.Email,
			}

			ctx := context.WithValue(req.Context(), UserContextKey{}, userInfo)
			req = req.WithContext(ctx)

			return next(w, req)
		}

		// If JWT validation fails, try OIDC token validation
		userInfo, err := m.authService.ValidateOIDCToken(req.Context(), token)
		if err == nil {
			// OIDC token is valid, add customer info to context
			customerInfo := &CustomerInfo{
				ID:    userInfo.Subject,
				Email: userInfo.Email,
			}

			ctx := context.WithValue(req.Context(), CustomerContextKey{}, customerInfo)
			req = req.WithContext(ctx)
		}

		// Continue regardless of token validation result
		return next(w, req)
	}
}

// GetCustomerFromContext extracts customer info from context
func GetCustomerFromContext(ctx context.Context) (*CustomerInfo, bool) {
	customer, ok := ctx.Value(CustomerContextKey{}).(*CustomerInfo)
	return customer, ok
}

// GetUserFromContext extracts user info from context
func GetUserFromContext(ctx context.Context) (*UserInfo, bool) {
	user, ok := ctx.Value(UserContextKey{}).(*UserInfo)
	return user, ok
}
