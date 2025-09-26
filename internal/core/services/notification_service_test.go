package services

import (
	"context"
	"testing"

	"silbackendassessment/internal/core/ports"
	"silbackendassessment/internal/testutils"
)

func TestNotificationService_SendEmail(t *testing.T) {
	mockEmailClient := &testutils.MockEmailClient{}
	mockSMSClient := &testutils.MockSMSClient{}
	service := NewNotificationService(mockEmailClient, mockSMSClient)

	ctx := context.Background()

	t.Run("Send email successfully", func(t *testing.T) {
		err := service.SendEmail(ctx, "test@example.com", "Test Subject", "Test Body", "Test HTML")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(mockEmailClient.SentEmails) != 1 {
			t.Errorf("Expected 1 email sent, got: %d", len(mockEmailClient.SentEmails))
		}

		email := mockEmailClient.SentEmails[0]
		if email.To != "test@example.com" {
			t.Errorf("Expected recipient to be test@example.com, got: %s", email.To)
		}
		if email.Subject != "Test Subject" {
			t.Errorf("Expected subject to be 'Test Subject', got: %s", email.Subject)
		}
		if email.Body != "Test Body" {
			t.Errorf("Expected body to be 'Test Body', got: %s", email.Body)
		}
		if email.HTMLBody != "Test HTML" {
			t.Errorf("Expected HTML body to be 'Test HTML', got: %s", email.HTMLBody)
		}
	})

	t.Run("Send email without HTML body", func(t *testing.T) {
		mockEmailClient.SentEmails = nil // Reset

		err := service.SendEmail(ctx, "test@example.com", "Test Subject", "Test Body")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(mockEmailClient.SentEmails) != 1 {
			t.Errorf("Expected 1 email sent, got: %d", len(mockEmailClient.SentEmails))
		}

		email := mockEmailClient.SentEmails[0]
		if email.HTMLBody != "" {
			t.Errorf("Expected empty HTML body, got: %s", email.HTMLBody)
		}
	})
}

func TestNotificationService_SendSMS(t *testing.T) {
	mockEmailClient := &testutils.MockEmailClient{}
	mockSMSClient := &testutils.MockSMSClient{}
	service := NewNotificationService(mockEmailClient, mockSMSClient)

	ctx := context.Background()

	t.Run("Send SMS successfully", func(t *testing.T) {
		err := service.SendSMS(ctx, "+1234567890", "Test message")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(mockSMSClient.SentSMS) != 1 {
			t.Errorf("Expected 1 SMS sent, got: %d", len(mockSMSClient.SentSMS))
		}

		sms := mockSMSClient.SentSMS[0]
		if sms.PhoneNumber != "+1234567890" {
			t.Errorf("Expected phone number to be +1234567890, got: %s", sms.PhoneNumber)
		}
		if sms.Message != "Test message" {
			t.Errorf("Expected message to be 'Test message', got: %s", sms.Message)
		}
	})
}

func TestNotificationService_SendNotification(t *testing.T) {
	mockEmailClient := &testutils.MockEmailClient{}
	mockSMSClient := &testutils.MockSMSClient{}
	service := NewNotificationService(mockEmailClient, mockSMSClient)

	ctx := context.Background()

	t.Run("Send email notification", func(t *testing.T) {
		req := &ports.NotificationRequest{
			Type:     ports.EmailNotification,
			To:       "test@example.com",
			Subject:  "Test Subject",
			Body:     "Test Body",
			HTMLBody: "Test HTML",
		}

		err := service.SendNotification(ctx, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(mockEmailClient.SentEmails) != 1 {
			t.Errorf("Expected 1 email sent, got: %d", len(mockEmailClient.SentEmails))
		}
	})

	t.Run("Send SMS notification", func(t *testing.T) {
		mockSMSClient.SentSMS = nil // Reset

		req := &ports.NotificationRequest{
			Type:        ports.SMSNotification,
			PhoneNumber: "+1234567890",
			Message:     "Test message",
		}

		err := service.SendNotification(ctx, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(mockSMSClient.SentSMS) != 1 {
			t.Errorf("Expected 1 SMS sent, got: %d", len(mockSMSClient.SentSMS))
		}
	})

	t.Run("Nil notification request", func(t *testing.T) {
		err := service.SendNotification(ctx, nil)

		if err == nil {
			t.Error("Expected error for nil request")
		}
		if err.Error() != "notification request cannot be nil" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})
}

func TestNotificationService_SendBulkEmail(t *testing.T) {
	mockEmailClient := &testutils.MockEmailClient{}
	mockSMSClient := &testutils.MockSMSClient{}
	service := NewNotificationService(mockEmailClient, mockSMSClient)

	ctx := context.Background()

	t.Run("Send bulk email successfully", func(t *testing.T) {
		recipients := []string{"test1@example.com", "test2@example.com", "test3@example.com"}

		err := service.SendBulkEmail(ctx, recipients, "Bulk Subject", "Bulk Body", "Bulk HTML")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(mockEmailClient.SentEmails) != 3 {
			t.Errorf("Expected 3 emails sent, got: %d", len(mockEmailClient.SentEmails))
		}

		// Verify all recipients received the email
		for i, recipient := range recipients {
			if mockEmailClient.SentEmails[i].To != recipient {
				t.Errorf("Expected recipient %d to be %s, got: %s", i, recipient, mockEmailClient.SentEmails[i].To)
			}
		}
	})

	t.Run("Empty recipients list", func(t *testing.T) {
		err := service.SendBulkEmail(ctx, []string{}, "Subject", "Body")

		if err == nil {
			t.Error("Expected error for empty recipients")
		}
		if err.Error() != "no recipients provided" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})
}

func TestNotificationService_SendBulkSMS(t *testing.T) {
	mockEmailClient := &testutils.MockEmailClient{}
	mockSMSClient := &testutils.MockSMSClient{}
	service := NewNotificationService(mockEmailClient, mockSMSClient)

	ctx := context.Background()

	t.Run("Send bulk SMS successfully", func(t *testing.T) {
		phoneNumbers := []string{"+1234567890", "+0987654321", "+1122334455"}

		err := service.SendBulkSMS(ctx, phoneNumbers, "Bulk message")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(mockSMSClient.SentSMS) != 3 {
			t.Errorf("Expected 3 SMS sent, got: %d", len(mockSMSClient.SentSMS))
		}

		// Verify all phone numbers received the SMS
		for i, phoneNumber := range phoneNumbers {
			if mockSMSClient.SentSMS[i].PhoneNumber != phoneNumber {
				t.Errorf("Expected phone number %d to be %s, got: %s", i, phoneNumber, mockSMSClient.SentSMS[i].PhoneNumber)
			}
		}
	})

	t.Run("Empty phone numbers list", func(t *testing.T) {
		err := service.SendBulkSMS(ctx, []string{}, "Message")

		if err == nil {
			t.Error("Expected error for empty phone numbers")
		}
		if err.Error() != "no phone numbers provided" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})
}

func TestNotificationService_Validation(t *testing.T) {
	mockEmailClient := &testutils.MockEmailClient{}
	mockSMSClient := &testutils.MockSMSClient{}
	service := NewNotificationService(mockEmailClient, mockSMSClient)

	t.Run("Validate email", func(t *testing.T) {
		result := service.ValidateEmail("test@example.com")
		if !result {
			t.Error("Expected valid email to return true")
		}

		result = service.ValidateEmail("invalid-email")
		if result {
			t.Error("Expected invalid email to return false")
		}
	})

	t.Run("Validate phone number", func(t *testing.T) {
		result := service.ValidatePhoneNumber("+1234567890")
		if !result {
			t.Error("Expected valid phone number to return true")
		}

		result = service.ValidatePhoneNumber("invalid")
		if result {
			t.Error("Expected invalid phone number to return false")
		}
	})
}
