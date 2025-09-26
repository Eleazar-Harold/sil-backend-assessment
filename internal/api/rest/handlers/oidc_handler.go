package handlers

import (
	"encoding/json"
	"net/http"

	"silbackendassessment/internal/core/ports"

	"github.com/uptrace/bunrouter"
)

// OIDCHandler handles OpenID Connect authentication requests
type OIDCHandler struct {
	authService ports.AuthService
}

// NewOIDCHandler creates a new OIDC handler
func NewOIDCHandler(authService ports.AuthService) *OIDCHandler {
	return &OIDCHandler{
		authService: authService,
	}
}

// AuthURLResponse represents the response for auth URL generation
type AuthURLResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

// OIDCLoginResponse represents the response for OIDC login
type OIDCLoginResponse struct {
	AccessToken  string                 `json:"access_token"`
	RefreshToken string                 `json:"refresh_token"`
	Customer     map[string]interface{} `json:"customer"`
	ExpiresAt    string                 `json:"expires_at"`
	IsNewUser    bool                   `json:"is_new_user"`
}

// GetAuthURL generates the OIDC authorization URL
func (h *OIDCHandler) GetAuthURL(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	authURL, state, err := h.authService.GetOIDCAuthURL(ctx)
	if err != nil {
		http.Error(w, "Failed to generate auth URL", http.StatusInternalServerError)
		return err
	}

	response := AuthURLResponse{
		AuthURL: authURL,
		State:   state,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// HandleCallback handles the OIDC callback
func (h *OIDCHandler) HandleCallback(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	// Parse query parameters
	query := req.URL.Query()
	code := query.Get("code")
	state := query.Get("state")
	errorParam := query.Get("error")

	// Check for OAuth errors
	if errorParam != "" {
		errorDescription := query.Get("error_description")
		http.Error(w, "OAuth error: "+errorParam+" - "+errorDescription, http.StatusBadRequest)
		return nil
	}

	// Validate required parameters
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return nil
	}

	if state == "" {
		http.Error(w, "Missing state parameter", http.StatusBadRequest)
		return nil
	}

	// Handle the OIDC callback
	loginResponse, err := h.authService.HandleOIDCCallback(ctx, code, state)
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return err
	}

	// Convert customer to map for JSON response
	customerMap := map[string]interface{}{
		"id":         loginResponse.Customer.ID,
		"first_name": loginResponse.Customer.FirstName,
		"last_name":  loginResponse.Customer.LastName,
		"email":      loginResponse.Customer.Email,
		"phone":      loginResponse.Customer.Phone,
		"address":    loginResponse.Customer.Address,
		"city":       loginResponse.Customer.City,
		"state":      loginResponse.Customer.State,
		"zip_code":   loginResponse.Customer.ZipCode,
		"country":    loginResponse.Customer.Country,
		"created_at": loginResponse.Customer.CreatedAt,
		"updated_at": loginResponse.Customer.UpdatedAt,
	}

	response := OIDCLoginResponse{
		AccessToken:  loginResponse.AccessToken,
		RefreshToken: loginResponse.RefreshToken,
		Customer:     customerMap,
		ExpiresAt:    loginResponse.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		IsNewUser:    loginResponse.IsNewUser,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// ValidateToken validates an OIDC token
func (h *OIDCHandler) ValidateToken(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

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

	// Validate the token
	userInfo, err := h.authService.ValidateOIDCToken(ctx, token)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return err
	}

	// Return user info
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(userInfo)
}

// Logout handles OIDC logout (redirects to provider logout)
func (h *OIDCHandler) Logout(w http.ResponseWriter, req bunrouter.Request) error {
	// In a real implementation, you would:
	// 1. Invalidate the session/token
	// 2. Redirect to the OIDC provider's logout endpoint
	// 3. Clear any local cookies/sessions

	// For now, just return a success response
	response := map[string]string{
		"message": "Logged out successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registers OIDC routes
func (h *OIDCHandler) RegisterRoutes(router *bunrouter.Router) {
	router.GET("/auth/oidc/login", h.GetAuthURL)
	router.GET("/auth/oidc/callback", h.HandleCallback)
	router.GET("/auth/oidc/validate", h.ValidateToken)
	router.POST("/auth/oidc/logout", h.Logout)
}
