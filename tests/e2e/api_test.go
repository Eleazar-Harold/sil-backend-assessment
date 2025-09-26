package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"silbackendassessment/internal/testutils"
)

const (
	baseURL = "http://localhost:8080"
)

type APITest struct {
	client *http.Client
	db     *testutils.TestDB
}

func NewAPITest(t *testing.T) *APITest {
	// Initialize test database
	db := testutils.NewTestDB(t)

	// Setup test schema
	if err := db.SetupTestSchema(); err != nil {
		t.Fatalf("Failed to setup test schema: %v", err)
	}

	// Cleanup before test
	if err := db.Cleanup(); err != nil {
		t.Fatalf("Failed to cleanup test database: %v", err)
	}

	return &APITest{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		db: db,
	}
}

func (a *APITest) Close() {
	a.db.Close()
}

func (a *APITest) makeRequest(method, endpoint string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return a.client.Do(req)
}

func TestHealthEndpoint(t *testing.T) {
	api := NewAPITest(t)
	defer api.Close()

	resp, err := api.makeRequest("GET", "/api/health", nil, nil)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got: %v", response["status"])
	}
}

func TestDocumentationEndpoints(t *testing.T) {
	api := NewAPITest(t)
	defer api.Close()

	tests := []struct {
		name     string
		endpoint string
		expected int
	}{
		{"Root endpoint", "/", http.StatusOK},
		{"Docs endpoint", "/docs", http.StatusOK},
		{"Swagger JSON", "/swagger.json", http.StatusOK},
		{"Redoc HTML", "/redoc.html", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := api.makeRequest("GET", tt.endpoint, nil, nil)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expected {
				t.Errorf("Expected status %d, got: %d", tt.expected, resp.StatusCode)
			}
		})
	}
}

func TestUserEndpoints(t *testing.T) {
	api := NewAPITest(t)
	defer api.Close()

	t.Run("Create user without authentication", func(t *testing.T) {
		userData := map[string]string{
			"name":  "Test User",
			"email": "test@example.com",
		}

		resp, err := api.makeRequest("POST", "/api/users", userData, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should require authentication
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got: %d", resp.StatusCode)
		}
	})

	t.Run("Get users without authentication", func(t *testing.T) {
		resp, err := api.makeRequest("GET", "/api/users", nil, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should require authentication
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got: %d", resp.StatusCode)
		}
	})
}

func TestCustomerEndpoints(t *testing.T) {
	api := NewAPITest(t)
	defer api.Close()

	t.Run("Create customer without authentication", func(t *testing.T) {
		customerData := map[string]interface{}{
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john.doe@example.com",
			"phone":      "+1234567890",
		}

		resp, err := api.makeRequest("POST", "/api/customers", customerData, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should require authentication
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got: %d", resp.StatusCode)
		}
	})
}

func TestNotificationEndpoints(t *testing.T) {
	api := NewAPITest(t)
	defer api.Close()

	t.Run("Send email without authentication", func(t *testing.T) {
		emailData := map[string]string{
			"to":      "test@example.com",
			"subject": "Test Subject",
			"body":    "Test Body",
		}

		resp, err := api.makeRequest("POST", "/api/notifications/email", emailData, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should require authentication
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got: %d", resp.StatusCode)
		}
	})

	t.Run("Send SMS without authentication", func(t *testing.T) {
		smsData := map[string]string{
			"phone_number": "+1234567890",
			"message":      "Test message",
		}

		resp, err := api.makeRequest("POST", "/api/notifications/sms", smsData, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should require authentication
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got: %d", resp.StatusCode)
		}
	})

	t.Run("Send bulk email without authentication", func(t *testing.T) {
		bulkEmailData := map[string]interface{}{
			"recipients": []string{"test1@example.com", "test2@example.com"},
			"subject":    "Bulk Test Subject",
			"body":       "Bulk Test Body",
		}

		resp, err := api.makeRequest("POST", "/api/notifications/bulk/email", bulkEmailData, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should require authentication
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got: %d", resp.StatusCode)
		}
	})
}

func TestGraphQLEndpoints(t *testing.T) {
	api := NewAPITest(t)
	defer api.Close()

	t.Run("GraphQL query without authentication", func(t *testing.T) {
		queryData := map[string]string{
			"query": "{ __schema { queryType { fields { name } } } }",
		}

		resp, err := api.makeRequest("POST", "/graphql", queryData, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should require authentication
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got: %d", resp.StatusCode)
		}
	})

	t.Run("GraphQL playground", func(t *testing.T) {
		resp, err := api.makeRequest("GET", "/graphql/playground", nil, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got: %d", resp.StatusCode)
		}
	})
}

func TestOIDCEndpoints(t *testing.T) {
	api := NewAPITest(t)
	defer api.Close()

	t.Run("OIDC login endpoint", func(t *testing.T) {
		resp, err := api.makeRequest("GET", "/auth/oidc/login", nil, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got: %d", resp.StatusCode)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if _, exists := response["auth_url"]; !exists {
			t.Error("Expected response to contain auth_url")
		}
	})

	t.Run("OIDC callback without parameters", func(t *testing.T) {
		resp, err := api.makeRequest("GET", "/auth/oidc/callback", nil, nil)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Should return error for missing parameters
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got: %d", resp.StatusCode)
		}
	})
}

func TestInvalidEndpoints(t *testing.T) {
	api := NewAPITest(t)
	defer api.Close()

	tests := []struct {
		name     string
		endpoint string
		method   string
		expected int
	}{
		{"Non-existent endpoint", "/api/nonexistent", "GET", http.StatusNotFound},
		{"Invalid method on valid endpoint", "/api/users", "PATCH", http.StatusMethodNotAllowed},
		{"GraphQL with invalid method", "/graphql", "GET", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := api.makeRequest(tt.method, tt.endpoint, nil, nil)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expected {
				t.Errorf("Expected status %d, got: %d", tt.expected, resp.StatusCode)
			}
		})
	}
}

func TestCORSHeaders(t *testing.T) {
	api := NewAPITest(t)
	defer api.Close()

	t.Run("CORS preflight request", func(t *testing.T) {
		req, err := http.NewRequest("OPTIONS", baseURL+"/api/users", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")

		resp, err := api.client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got: %d", resp.StatusCode)
		}

		// Check CORS headers
		corsHeaders := []string{
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods",
			"Access-Control-Allow-Headers",
		}

		for _, header := range corsHeaders {
			if resp.Header.Get(header) == "" {
				t.Errorf("Expected CORS header %s to be present", header)
			}
		}
	})
}
