package notifications

import (
	"context"
	"strings"
	"testing"
)

func TestEmailClient_ValidateEmail(t *testing.T) {
	client := NewEmailClient(&SMTPConfig{
		Host: "smtp.gmail.com",
		Port: 587,
	})

	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Valid email", "test@example.com", true},
		{"Valid email with subdomain", "user@mail.example.com", true},
		{"Valid email with plus", "user+tag@example.com", true},
		{"Valid email with dot", "user.name@example.com", true},
		{"Invalid email - no @", "testexample.com", false},
		{"Invalid email - no domain", "test@", false},
		{"Invalid email - no local", "@example.com", false},
		{"Invalid email - empty", "", false},
		{"Invalid email - spaces", "test @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.ValidateEmail(tt.email)
			if result != tt.expected {
				t.Errorf("ValidateEmail(%s) = %v, expected %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestEmailClient_BuildMessage(t *testing.T) {
	client := NewEmailClient(&SMTPConfig{
		Host: "smtp.gmail.com",
		Port: 587,
		From: "sender@example.com",
	})

	t.Run("Plain text message", func(t *testing.T) {
		message := client.buildMessage("recipient@example.com", "Test Subject", "Test Body")

		if message == "" {
			t.Error("Expected non-empty message")
		}

		// Check for required headers
		requiredHeaders := []string{"From:", "To:", "Subject:", "Date:", "MIME-Version:"}
		for _, header := range requiredHeaders {
			if !strings.Contains(message, header) {
				t.Errorf("Expected message to contain header %s", header)
			}
		}

		if !strings.Contains(message, "Test Subject") {
			t.Error("Expected message to contain subject")
		}

		if !strings.Contains(message, "Test Body") {
			t.Error("Expected message to contain body")
		}
	})

	t.Run("HTML message", func(t *testing.T) {
		htmlBody := "<h1>Test</h1><p>HTML body</p>"
		message := client.buildMessage("recipient@example.com", "Test Subject", "Test Body", htmlBody)

		if message == "" {
			t.Error("Expected non-empty message")
		}

		if !strings.Contains(message, "multipart/alternative") {
			t.Error("Expected multipart message for HTML content")
		}

		if !strings.Contains(message, htmlBody) {
			t.Error("Expected message to contain HTML body")
		}

		if !strings.Contains(message, "Test Body") {
			t.Error("Expected message to contain plain text body")
		}
	})
}

func TestEmailClient_SendEmail_Validation(t *testing.T) {
	client := NewEmailClient(&SMTPConfig{
		Host: "smtp.gmail.com",
		Port: 587,
		From: "invalid-email",
	})

	ctx := context.Background()

	t.Run("Invalid recipient email", func(t *testing.T) {
		err := client.SendEmail(ctx, "invalid-email", "Subject", "Body")
		if err == nil {
			t.Error("Expected error for invalid recipient email")
		}
		if !strings.Contains(err.Error(), "invalid recipient email") {
			t.Errorf("Expected error about invalid recipient email, got: %v", err)
		}
	})

	t.Run("Invalid sender email", func(t *testing.T) {
		err := client.SendEmail(ctx, "test@example.com", "Subject", "Body")
		if err == nil {
			t.Error("Expected error for invalid sender email")
		}
		if !strings.Contains(err.Error(), "invalid sender email") {
			t.Errorf("Expected error about invalid sender email, got: %v", err)
		}
	})
}

func TestEmailClient_SendBulkEmail_Validation(t *testing.T) {
	client := NewEmailClient(&SMTPConfig{
		Host: "smtp.gmail.com",
		Port: 587,
		From: "sender@example.com",
	})

	ctx := context.Background()

	t.Run("Empty recipients", func(t *testing.T) {
		err := client.SendBulkEmail(ctx, []string{}, "Subject", "Body")
		if err == nil {
			t.Error("Expected error for empty recipients")
		}
		if !strings.Contains(err.Error(), "no recipients provided") {
			t.Errorf("Expected error about no recipients, got: %v", err)
		}
	})

	t.Run("Invalid recipient email", func(t *testing.T) {
		err := client.SendBulkEmail(ctx, []string{"invalid-email"}, "Subject", "Body")
		if err == nil {
			t.Error("Expected error for invalid recipient email")
		}
		if !strings.Contains(err.Error(), "invalid email address") {
			t.Errorf("Expected error about invalid email address, got: %v", err)
		}
	})
}

func TestEmailClient_SendEmail_EdgeCases(t *testing.T) {
	client := NewEmailClient(&SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: "test@example.com",
		Password: "password",
		From:     "sender@example.com",
		TLS:      true,
	})

	t.Run("Send email with empty recipient", func(t *testing.T) {
		err := client.SendEmail(context.Background(), "", "Test Subject", "Test Body")
		if err == nil {
			t.Error("Expected error for empty recipient")
		}
	})

	t.Run("Send email with invalid recipient", func(t *testing.T) {
		err := client.SendEmail(context.Background(), "invalid-email", "Test Subject", "Test Body")
		if err == nil {
			t.Error("Expected error for invalid recipient")
		}
	})

	t.Run("Send email with empty subject", func(t *testing.T) {
		err := client.SendEmail(context.Background(), "test@example.com", "", "Test Body")
		if err == nil {
			t.Error("Expected error for empty subject")
		}
	})

	t.Run("Send email with empty body", func(t *testing.T) {
		err := client.SendEmail(context.Background(), "test@example.com", "Test Subject", "")
		if err == nil {
			t.Error("Expected error for empty body")
		}
	})

	t.Run("Send email with very long subject", func(t *testing.T) {
		longSubject := strings.Repeat("A", 1000)
		err := client.SendEmail(context.Background(), "test@example.com", longSubject, "Test Body")
		if err == nil {
			t.Error("Expected error for very long subject")
		}
	})

	t.Run("Send email with very long body", func(t *testing.T) {
		longBody := strings.Repeat("B", 100000)
		err := client.SendEmail(context.Background(), "test@example.com", "Test Subject", longBody)
		if err == nil {
			t.Error("Expected error for very long body")
		}
	})
}

func TestEmailClient_SendBulkEmail_EdgeCases(t *testing.T) {
	client := NewEmailClient(&SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: "test@example.com",
		Password: "password",
		From:     "sender@example.com",
		TLS:      true,
	})

	t.Run("Send bulk email with empty recipients", func(t *testing.T) {
		err := client.SendBulkEmail(context.Background(), []string{}, "Test Subject", "Test Body")
		if err == nil {
			t.Error("Expected error for empty recipients")
		}
	})

	t.Run("Send bulk email with invalid recipients", func(t *testing.T) {
		recipients := []string{"invalid-email", "another-invalid"}
		err := client.SendBulkEmail(context.Background(), recipients, "Test Subject", "Test Body")
		if err == nil {
			t.Error("Expected error for invalid recipients")
		}
	})

	t.Run("Send bulk email with mixed valid and invalid recipients", func(t *testing.T) {
		recipients := []string{"valid@example.com", "invalid-email", "another@example.com"}
		err := client.SendBulkEmail(context.Background(), recipients, "Test Subject", "Test Body")
		if err == nil {
			t.Error("Expected error for mixed valid and invalid recipients")
		}
	})

	t.Run("Send bulk email with duplicate recipients", func(t *testing.T) {
		recipients := []string{"test@example.com", "test@example.com", "test@example.com"}
		err := client.SendBulkEmail(context.Background(), recipients, "Test Subject", "Test Body")
		if err == nil {
			t.Error("Expected error for duplicate recipients")
		}
	})

	t.Run("Send bulk email with too many recipients", func(t *testing.T) {
		recipients := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			recipients[i] = "test@example.com"
		}
		err := client.SendBulkEmail(context.Background(), recipients, "Test Subject", "Test Body")
		if err == nil {
			t.Error("Expected error for too many recipients")
		}
	})
}

func TestEmailClient_BuildMessage_EdgeCases(t *testing.T) {
	client := NewEmailClient(&SMTPConfig{
		Host: "smtp.gmail.com",
		Port: 587,
		From: "sender@example.com",
	})

	t.Run("Build message with special characters", func(t *testing.T) {
		message := client.buildMessage("test@example.com", "Subject with émojis 🎉", "Body with special chars: àáâãäå")

		if !strings.Contains(message, "Subject with émojis 🎉") {
			t.Error("Message should contain special characters in subject")
		}
		if !strings.Contains(message, "Body with special chars: àáâãäå") {
			t.Error("Message should contain special characters in body")
		}
	})

	t.Run("Build message with HTML content", func(t *testing.T) {
		htmlBody := "<html><body><h1>Test</h1><p>This is a test email</p></body></html>"
		message := client.buildMessage("test@example.com", "HTML Subject", htmlBody)

		if !strings.Contains(message, "HTML Subject") {
			t.Error("Message should contain HTML subject")
		}
		if !strings.Contains(message, htmlBody) {
			t.Error("Message should contain HTML body")
		}
	})

	t.Run("Build message with long content", func(t *testing.T) {
		longSubject := strings.Repeat("A", 200)
		longBody := strings.Repeat("B", 1000)
		message := client.buildMessage("test@example.com", longSubject, longBody)

		if !strings.Contains(message, longSubject) {
			t.Error("Message should contain long subject")
		}
		if !strings.Contains(message, longBody) {
			t.Error("Message should contain long body")
		}
	})
}
