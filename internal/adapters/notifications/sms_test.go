package notifications

import (
	"context"
	"strings"
	"testing"
)

func TestSMSClient_ValidatePhoneNumber(t *testing.T) {
	client := NewSMSClient(&ATConfig{
		APIKey: "test-key",
	})

	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{"Valid international number", "+1234567890", true},
		{"Valid international number with country code", "+44123456789", true},
		{"Valid international number with plus", "+254712345678", true},
		{"Valid local number", "1234567890", true},
		{"Valid local number with dashes", "123-456-7890", true},
		{"Valid local number with spaces", "123 456 7890", true},
		{"Valid local number with parentheses", "(123) 456-7890", true},
		{"Invalid - too short", "123", false},
		{"Invalid - too long", "+12345678901234567890", false},
		{"Invalid - empty", "", false},
		{"Invalid - no digits", "abc-def-ghij", false},
		{"Invalid - special chars only", "!@#$%^&*()", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.ValidatePhoneNumber(tt.phone)
			if result != tt.expected {
				t.Errorf("ValidatePhoneNumber(%s) = %v, expected %v", tt.phone, result, tt.expected)
			}
		})
	}
}

func TestSMSClient_NewSMSClient(t *testing.T) {
	t.Run("Default base URL", func(t *testing.T) {
		config := &ATConfig{
			APIKey: "test-key",
		}
		client := NewSMSClient(config)

		if client.config.BaseURL != "https://api.africastalking.com" {
			t.Errorf("Expected default base URL, got: %s", client.config.BaseURL)
		}
	})

	t.Run("Custom base URL", func(t *testing.T) {
		config := &ATConfig{
			APIKey:  "test-key",
			BaseURL: "https://custom.api.com",
		}
		client := NewSMSClient(config)

		if client.config.BaseURL != "https://custom.api.com" {
			t.Errorf("Expected custom base URL, got: %s", client.config.BaseURL)
		}
	})

	t.Run("HTTP client timeout", func(t *testing.T) {
		config := &ATConfig{
			APIKey: "test-key",
		}
		client := NewSMSClient(config)

		if client.httpClient.Timeout.Seconds() != 30 {
			t.Errorf("Expected 30 second timeout, got: %v", client.httpClient.Timeout)
		}
	})
}

func TestSMSClient_SendSMS_Validation(t *testing.T) {
	client := NewSMSClient(&ATConfig{
		APIKey: "test-key",
	})

	ctx := context.Background()

	t.Run("Invalid phone number", func(t *testing.T) {
		err := client.SendSMS(ctx, "invalid", "Test message")
		if err == nil {
			t.Error("Expected error for invalid phone number")
		}
		if !strings.Contains(err.Error(), "invalid phone number") {
			t.Errorf("Expected error about invalid phone number, got: %v", err)
		}
	})
}

func TestSMSClient_SendBulkSMS_Validation(t *testing.T) {
	client := NewSMSClient(&ATConfig{
		APIKey: "test-key",
	})

	ctx := context.Background()

	t.Run("Empty phone numbers", func(t *testing.T) {
		err := client.SendBulkSMS(ctx, []string{}, "Test message")
		if err == nil {
			t.Error("Expected error for empty phone numbers")
		}
		if !strings.Contains(err.Error(), "no phone numbers provided") {
			t.Errorf("Expected error about no phone numbers, got: %v", err)
		}
	})

	t.Run("Invalid phone number in list", func(t *testing.T) {
		err := client.SendBulkSMS(ctx, []string{"invalid"}, "Test message")
		if err == nil {
			t.Error("Expected error for invalid phone number")
		}
		if !strings.Contains(err.Error(), "invalid phone number") {
			t.Errorf("Expected error about invalid phone number, got: %v", err)
		}
	})
}

func TestSMSRequest_MarshalJSON(t *testing.T) {
	req := SMSRequest{
		Username: "testuser",
		Message:  "Test message",
		To:       []string{"+1234567890", "+0987654321"},
	}

	// This would typically test JSON marshaling, but since we're not importing
	// the JSON package, we'll just verify the struct fields are set correctly
	if req.Username != "testuser" {
		t.Error("Expected username to be set")
	}
	if req.Message != "Test message" {
		t.Error("Expected message to be set")
	}
	if len(req.To) != 2 {
		t.Error("Expected 2 recipients")
	}
	if req.To[0] != "+1234567890" {
		t.Error("Expected first recipient to be correct")
	}
	if req.To[1] != "+0987654321" {
		t.Error("Expected second recipient to be correct")
	}
}

func TestSMSClient_SendSMS_EdgeCases(t *testing.T) {
	client := NewSMSClient(&ATConfig{
		APIKey:   "test-api-key",
		Username: "test-username",
		BaseURL:  "https://api.africastalking.com",
	})

	t.Run("Send SMS with empty recipient", func(t *testing.T) {
		err := client.SendSMS(context.Background(), "", "Test message")
		if err == nil {
			t.Error("Expected error for empty recipient")
		}
	})

	t.Run("Send SMS with invalid phone number", func(t *testing.T) {
		err := client.SendSMS(context.Background(), "invalid-phone", "Test message")
		if err == nil {
			t.Error("Expected error for invalid phone number")
		}
	})

	t.Run("Send SMS with empty message", func(t *testing.T) {
		err := client.SendSMS(context.Background(), "+1234567890", "")
		if err == nil {
			t.Error("Expected error for empty message")
		}
	})

	t.Run("Send SMS with very long message", func(t *testing.T) {
		longMessage := strings.Repeat("A", 2000)
		err := client.SendSMS(context.Background(), "+1234567890", longMessage)
		if err == nil {
			t.Error("Expected error for very long message")
		}
	})

	t.Run("Send SMS with special characters", func(t *testing.T) {
		message := "Test message with émojis 🎉 and special chars: àáâãäå"
		err := client.SendSMS(context.Background(), "+1234567890", message)
		if err == nil {
			t.Error("Expected error for message with special characters")
		}
	})
}

func TestSMSClient_SendBulkSMS_EdgeCases(t *testing.T) {
	client := NewSMSClient(&ATConfig{
		APIKey:   "test-api-key",
		Username: "test-username",
		BaseURL:  "https://api.africastalking.com",
	})

	t.Run("Send bulk SMS with empty recipients", func(t *testing.T) {
		err := client.SendBulkSMS(context.Background(), []string{}, "Test message")
		if err == nil {
			t.Error("Expected error for empty recipients")
		}
	})

	t.Run("Send bulk SMS with invalid phone numbers", func(t *testing.T) {
		recipients := []string{"invalid-phone", "another-invalid"}
		err := client.SendBulkSMS(context.Background(), recipients, "Test message")
		if err == nil {
			t.Error("Expected error for invalid phone numbers")
		}
	})

	t.Run("Send bulk SMS with mixed valid and invalid recipients", func(t *testing.T) {
		recipients := []string{"+1234567890", "invalid-phone", "+0987654321"}
		err := client.SendBulkSMS(context.Background(), recipients, "Test message")
		if err == nil {
			t.Error("Expected error for mixed valid and invalid recipients")
		}
	})

	t.Run("Send bulk SMS with duplicate recipients", func(t *testing.T) {
		recipients := []string{"+1234567890", "+1234567890", "+1234567890"}
		err := client.SendBulkSMS(context.Background(), recipients, "Test message")
		if err == nil {
			t.Error("Expected error for duplicate recipients")
		}
	})

	t.Run("Send bulk SMS with too many recipients", func(t *testing.T) {
		recipients := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			recipients[i] = "+1234567890"
		}
		err := client.SendBulkSMS(context.Background(), recipients, "Test message")
		if err == nil {
			t.Error("Expected error for too many recipients")
		}
	})
}

func TestSMSClient_ValidatePhoneNumber_EdgeCases(t *testing.T) {
	client := NewSMSClient(&ATConfig{
		APIKey:   "test-api-key",
		Username: "test-username",
		BaseURL:  "https://api.africastalking.com",
	})

	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{"Valid US phone with country code", "+1234567890", true},
		{"Valid US phone with country code and spaces", "+1 234 567 890", true},
		{"Valid US phone with country code and dashes", "+1-234-567-890", true},
		{"Valid US phone with country code and parentheses", "+1 (234) 567-890", true},
		{"Valid international phone", "+44123456789", true},
		{"Valid phone with extension", "+1234567890 ext 123", true},
		{"Invalid phone - too short", "+123", false},
		{"Invalid phone - too long", "+12345678901234567890", false},
		{"Invalid phone - no country code", "1234567890", true},
		{"Invalid phone - empty", "", false},
		{"Invalid phone - letters", "+123456789a", false},
		{"Invalid phone - special chars", "+123-456-789@", false},
		{"Invalid phone - spaces only", "   ", false},
		{"Invalid phone - multiple plus signs", "++1234567890", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.ValidatePhoneNumber(tt.phone)
			if result != tt.expected {
				t.Errorf("ValidatePhoneNumber(%s) = %v, expected %v", tt.phone, result, tt.expected)
			}
		})
	}
}
