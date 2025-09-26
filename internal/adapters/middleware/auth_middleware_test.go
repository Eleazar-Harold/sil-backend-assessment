package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"silbackendassessment/internal/adapters/auth"
	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/oidc"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
	"github.com/uptrace/bunrouter"
)

// MockAuthService for testing
type MockAuthService struct {
	ValidateTokenFunc      func(token string) (*auth.Claims, error)
	ValidateOIDCTokenFunc  func(ctx context.Context, token string) (*oidc.OIDCUserInfo, error)
	GetOIDCAuthURLFunc     func(ctx context.Context) (string, string, error)
	HandleOIDCCallbackFunc func(ctx context.Context, code, state string) (*ports.OIDCLoginResponse, error)
	LoginFunc              func(ctx context.Context, req *ports.LoginRequest) (*ports.LoginResponse, error)
	RegisterFunc           func(ctx context.Context, req *ports.RegisterRequest) (*ports.LoginResponse, error)
	RefreshTokenFunc       func(ctx context.Context, refreshToken string) (*ports.LoginResponse, error)
}

func (m *MockAuthService) ValidateToken(token string) (*auth.Claims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(token)
	}
	return &auth.Claims{
		UserID: "user-123",
		Email:  "test@example.com",
	}, nil
}

func (m *MockAuthService) ValidateOIDCToken(ctx context.Context, token string) (*oidc.OIDCUserInfo, error) {
	if m.ValidateOIDCTokenFunc != nil {
		return m.ValidateOIDCTokenFunc(ctx, token)
	}
	return &oidc.OIDCUserInfo{
		Subject: "user-123",
		Email:   "test@example.com",
	}, nil
}

func (m *MockAuthService) GetOIDCAuthURL(ctx context.Context) (string, string, error) {
	if m.GetOIDCAuthURLFunc != nil {
		return m.GetOIDCAuthURLFunc(ctx)
	}
	return "http://example.com/auth", "state-123", nil
}

func (m *MockAuthService) HandleOIDCCallback(ctx context.Context, code, state string) (*ports.OIDCLoginResponse, error) {
	if m.HandleOIDCCallbackFunc != nil {
		return m.HandleOIDCCallbackFunc(ctx, code, state)
	}
	return &ports.OIDCLoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Customer: &domain.Customer{
			ID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Email: "customer@example.com",
		},
		ExpiresAt: time.Now().Add(time.Hour),
		IsNewUser: false,
	}, nil
}

func (m *MockAuthService) Login(ctx context.Context, req *ports.LoginRequest) (*ports.LoginResponse, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, req)
	}
	return &ports.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		User: &domain.User{
			ID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			Email: "test@example.com",
		},
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func (m *MockAuthService) Register(ctx context.Context, req *ports.RegisterRequest) (*ports.LoginResponse, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, req)
	}
	return &ports.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		User: &domain.User{
			ID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			Email: "test@example.com",
		},
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*ports.LoginResponse, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, refreshToken)
	}
	return &ports.LoginResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		User: &domain.User{
			ID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			Email: "test@example.com",
		},
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func TestAuthMiddleware_NewAuthMiddleware(t *testing.T) {
	mockAuthService := &MockAuthService{}
	middleware := NewAuthMiddleware(mockAuthService)

	if middleware == nil {
		t.Error("Expected middleware to be non-nil")
	}
	if middleware.authService == nil {
		t.Error("Expected authService to be set")
	}
}

func TestAuthMiddleware_RequireAuth(t *testing.T) {
	t.Run("Valid JWT token", func(t *testing.T) {
		mockAuthService := &MockAuthService{
			ValidateTokenFunc: func(token string) (*auth.Claims, error) {
				if token == "valid-token" {
					return &auth.Claims{
						UserID: "user-123",
						Email:  "test@example.com",
					}, nil
				}
				return nil, errors.New("invalid token")
			},
		}

		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.RequireAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			// Check if user info is in context
			user, ok := GetUserFromContext(req.Context())
			if !ok {
				t.Error("Expected user info to be in context")
			}
			if user.ID != "user-123" {
				t.Errorf("Expected user ID to be 'user-123', got: %s", user.ID)
			}
			if user.Email != "test@example.com" {
				t.Errorf("Expected user email to be 'test@example.com', got: %s", user.Email)
			}
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Missing Authorization header", func(t *testing.T) {
		mockAuthService := &MockAuthService{}
		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.RequireAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			t.Error("Handler should not be called")
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got: %d", http.StatusUnauthorized, w.Code)
		}
		if w.Body.String() != "Missing Authorization header\n" {
			t.Errorf("Expected 'Missing Authorization header', got: %s", w.Body.String())
		}
	})

	t.Run("Invalid Authorization header format", func(t *testing.T) {
		mockAuthService := &MockAuthService{}
		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.RequireAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			t.Error("Handler should not be called")
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got: %d", http.StatusUnauthorized, w.Code)
		}
		if w.Body.String() != "Invalid Authorization header format\n" {
			t.Errorf("Expected 'Invalid Authorization header format', got: %s", w.Body.String())
		}
	})

	t.Run("Invalid JWT token", func(t *testing.T) {
		mockAuthService := &MockAuthService{
			ValidateTokenFunc: func(token string) (*auth.Claims, error) {
				return nil, errors.New("invalid token")
			},
		}

		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.RequireAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			t.Error("Handler should not be called")
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err == nil {
			t.Error("Expected error for invalid token")
		}

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got: %d", http.StatusUnauthorized, w.Code)
		}
		if w.Body.String() != "Invalid token: invalid token\n" {
			t.Errorf("Expected 'Invalid token: invalid token', got: %s", w.Body.String())
		}
	})
}

func TestAuthMiddleware_RequireCustomerAuth(t *testing.T) {
	t.Run("Valid JWT token", func(t *testing.T) {
		mockAuthService := &MockAuthService{
			ValidateTokenFunc: func(token string) (*auth.Claims, error) {
				if token == "valid-token" {
					return &auth.Claims{
						UserID: "user-123",
						Email:  "test@example.com",
					}, nil
				}
				return nil, errors.New("invalid token")
			},
		}

		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.RequireCustomerAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			// Check if customer info is in context
			customer, ok := GetCustomerFromContext(req.Context())
			if !ok {
				t.Error("Expected customer info to be in context")
			}
			if customer.ID != "user-123" {
				t.Errorf("Expected customer ID to be 'user-123', got: %s", customer.ID)
			}
			if customer.Email != "test@example.com" {
				t.Errorf("Expected customer email to be 'test@example.com', got: %s", customer.Email)
			}
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Valid OIDC token", func(t *testing.T) {
		mockAuthService := &MockAuthService{
			ValidateTokenFunc: func(token string) (*auth.Claims, error) {
				return nil, errors.New("invalid JWT token")
			},
			ValidateOIDCTokenFunc: func(ctx context.Context, token string) (*oidc.OIDCUserInfo, error) {
				if token == "valid-oidc-token" {
					return &oidc.OIDCUserInfo{
						Subject: "user-456",
						Email:   "oidc@example.com",
					}, nil
				}
				return nil, errors.New("invalid OIDC token")
			},
		}

		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.RequireCustomerAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			// Check if customer info is in context
			customer, ok := GetCustomerFromContext(req.Context())
			if !ok {
				t.Error("Expected customer info to be in context")
			}
			if customer.ID != "user-456" {
				t.Errorf("Expected customer ID to be 'user-456', got: %s", customer.ID)
			}
			if customer.Email != "oidc@example.com" {
				t.Errorf("Expected customer email to be 'oidc@example.com', got: %s", customer.Email)
			}
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-oidc-token")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Invalid both JWT and OIDC tokens", func(t *testing.T) {
		mockAuthService := &MockAuthService{
			ValidateTokenFunc: func(token string) (*auth.Claims, error) {
				return nil, errors.New("invalid JWT token")
			},
			ValidateOIDCTokenFunc: func(ctx context.Context, token string) (*oidc.OIDCUserInfo, error) {
				return nil, errors.New("invalid OIDC token")
			},
		}

		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.RequireCustomerAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			t.Error("Handler should not be called")
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err == nil {
			t.Error("Expected error for invalid token")
		}

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got: %d", http.StatusUnauthorized, w.Code)
		}
		if w.Body.String() != "Invalid token: invalid OIDC token\n" {
			t.Errorf("Expected 'Invalid token: invalid OIDC token', got: %s", w.Body.String())
		}
	})
}

func TestAuthMiddleware_RequireOIDCAuth(t *testing.T) {
	t.Run("Valid OIDC token", func(t *testing.T) {
		mockAuthService := &MockAuthService{
			ValidateOIDCTokenFunc: func(ctx context.Context, token string) (*oidc.OIDCUserInfo, error) {
				if token == "valid-oidc-token" {
					return &oidc.OIDCUserInfo{
						Subject: "user-789",
						Email:   "oidc@example.com",
					}, nil
				}
				return nil, errors.New("invalid OIDC token")
			},
		}

		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.RequireOIDCAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			// Check if customer info is in context
			customer, ok := GetCustomerFromContext(req.Context())
			if !ok {
				t.Error("Expected customer info to be in context")
			}
			if customer.ID != "user-789" {
				t.Errorf("Expected customer ID to be 'user-789', got: %s", customer.ID)
			}
			if customer.Email != "oidc@example.com" {
				t.Errorf("Expected customer email to be 'oidc@example.com', got: %s", customer.Email)
			}
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-oidc-token")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Invalid OIDC token", func(t *testing.T) {
		mockAuthService := &MockAuthService{
			ValidateOIDCTokenFunc: func(ctx context.Context, token string) (*oidc.OIDCUserInfo, error) {
				return nil, errors.New("invalid OIDC token")
			},
		}

		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.RequireOIDCAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			t.Error("Handler should not be called")
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-oidc-token")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err == nil {
			t.Error("Expected error for invalid token")
		}

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got: %d", http.StatusUnauthorized, w.Code)
		}
		if w.Body.String() != "Invalid OIDC token: invalid OIDC token\n" {
			t.Errorf("Expected 'Invalid OIDC token: invalid OIDC token', got: %s", w.Body.String())
		}
	})
}

func TestAuthMiddleware_OptionalAuth(t *testing.T) {
	t.Run("No Authorization header", func(t *testing.T) {
		mockAuthService := &MockAuthService{}
		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.OptionalAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			// Should continue without authentication
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Invalid Authorization header format", func(t *testing.T) {
		mockAuthService := &MockAuthService{}
		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.OptionalAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			// Should continue without authentication
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Valid JWT token", func(t *testing.T) {
		mockAuthService := &MockAuthService{
			ValidateTokenFunc: func(token string) (*auth.Claims, error) {
				if token == "valid-token" {
					return &auth.Claims{
						UserID: "user-123",
						Email:  "test@example.com",
					}, nil
				}
				return nil, errors.New("invalid token")
			},
		}

		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.OptionalAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			// Check if user info is in context
			user, ok := GetUserFromContext(req.Context())
			if !ok {
				t.Error("Expected user info to be in context")
			}
			if user.ID != "user-123" {
				t.Errorf("Expected user ID to be 'user-123', got: %s", user.ID)
			}
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Valid OIDC token", func(t *testing.T) {
		mockAuthService := &MockAuthService{
			ValidateTokenFunc: func(token string) (*auth.Claims, error) {
				return nil, errors.New("invalid JWT token")
			},
			ValidateOIDCTokenFunc: func(ctx context.Context, token string) (*oidc.OIDCUserInfo, error) {
				if token == "valid-oidc-token" {
					return &oidc.OIDCUserInfo{
						Subject: "user-456",
						Email:   "oidc@example.com",
					}, nil
				}
				return nil, errors.New("invalid OIDC token")
			},
		}

		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.OptionalAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			// Check if customer info is in context
			customer, ok := GetCustomerFromContext(req.Context())
			if !ok {
				t.Error("Expected customer info to be in context")
			}
			if customer.ID != "user-456" {
				t.Errorf("Expected customer ID to be 'user-456', got: %s", customer.ID)
			}
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-oidc-token")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		mockAuthService := &MockAuthService{
			ValidateTokenFunc: func(token string) (*auth.Claims, error) {
				return nil, errors.New("invalid JWT token")
			},
			ValidateOIDCTokenFunc: func(ctx context.Context, token string) (*oidc.OIDCUserInfo, error) {
				return nil, errors.New("invalid OIDC token")
			},
		}

		middleware := NewAuthMiddleware(mockAuthService)
		handler := middleware.OptionalAuth(func(w http.ResponseWriter, req bunrouter.Request) error {
			// Should continue even with invalid token
			return nil
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		bunReq := bunrouter.NewRequest(req)
		err := handler(w, bunReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})
}

func TestGetCustomerFromContext(t *testing.T) {
	t.Run("Customer info in context", func(t *testing.T) {
		customer := &CustomerInfo{
			ID:    "customer-123",
			Email: "customer@example.com",
		}
		ctx := context.WithValue(context.Background(), CustomerContextKey{}, customer)

		retrievedCustomer, ok := GetCustomerFromContext(ctx)
		if !ok {
			t.Error("Expected to find customer in context")
		}
		if retrievedCustomer.ID != "customer-123" {
			t.Errorf("Expected customer ID 'customer-123', got: %s", retrievedCustomer.ID)
		}
		if retrievedCustomer.Email != "customer@example.com" {
			t.Errorf("Expected customer email 'customer@example.com', got: %s", retrievedCustomer.Email)
		}
	})

	t.Run("No customer info in context", func(t *testing.T) {
		ctx := context.Background()

		_, ok := GetCustomerFromContext(ctx)
		if ok {
			t.Error("Expected not to find customer in context")
		}
	})
}

func TestGetUserFromContext(t *testing.T) {
	t.Run("User info in context", func(t *testing.T) {
		user := &UserInfo{
			ID:    "user-123",
			Email: "user@example.com",
		}
		ctx := context.WithValue(context.Background(), UserContextKey{}, user)

		retrievedUser, ok := GetUserFromContext(ctx)
		if !ok {
			t.Error("Expected to find user in context")
		}
		if retrievedUser.ID != "user-123" {
			t.Errorf("Expected user ID 'user-123', got: %s", retrievedUser.ID)
		}
		if retrievedUser.Email != "user@example.com" {
			t.Errorf("Expected user email 'user@example.com', got: %s", retrievedUser.Email)
		}
	})

	t.Run("No user info in context", func(t *testing.T) {
		ctx := context.Background()

		_, ok := GetUserFromContext(ctx)
		if ok {
			t.Error("Expected not to find user in context")
		}
	})
}
