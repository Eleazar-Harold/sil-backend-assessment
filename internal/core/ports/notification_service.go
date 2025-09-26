package ports

import "context"

// EmailClient defines the interface for sending emails
type EmailClient interface {
	SendEmail(ctx context.Context, to, subject, body string, htmlBody ...string) error
	SendBulkEmail(ctx context.Context, recipients []string, subject, body string, htmlBody ...string) error
	ValidateEmail(email string) bool
}

// SMSClient defines the interface for sending SMS messages
type SMSClient interface {
	SendSMS(ctx context.Context, phoneNumber, message string) error
	SendBulkSMS(ctx context.Context, phoneNumbers []string, message string) error
	ValidatePhoneNumber(phoneNumber string) bool
}

// NotificationType represents the type of notification
type NotificationType string

const (
	EmailNotification NotificationType = "email"
	SMSNotification   NotificationType = "sms"
)

// NotificationRequest represents a notification request
type NotificationRequest struct {
	Type        NotificationType
	To          string
	Subject     string // For email
	Body        string
	HTMLBody    string // For email (optional)
	PhoneNumber string // For SMS
	Message     string // For SMS
}

// NotificationService defines the interface for sending notifications
type NotificationService interface {
	// SendEmail sends an email notification
	SendEmail(ctx context.Context, to, subject, body string, htmlBody ...string) error

	// SendSMS sends an SMS notification
	SendSMS(ctx context.Context, phoneNumber, message string) error

	// SendNotification sends a notification based on the request type
	SendNotification(ctx context.Context, req *NotificationRequest) error

	// SendBulkEmail sends emails to multiple recipients
	SendBulkEmail(ctx context.Context, recipients []string, subject, body string, htmlBody ...string) error

	// SendBulkSMS sends SMS to multiple phone numbers
	SendBulkSMS(ctx context.Context, phoneNumbers []string, message string) error

	// ValidateEmail validates if an email address is properly formatted
	ValidateEmail(email string) bool

	// ValidatePhoneNumber validates if a phone number is properly formatted
	ValidatePhoneNumber(phoneNumber string) bool
}
