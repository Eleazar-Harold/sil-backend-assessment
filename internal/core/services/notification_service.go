package services

import (
	"context"
	"fmt"
	"log"

	"silbackendassessment/internal/core/ports"
)

// NotificationService implements the notification service interface
type NotificationService struct {
	emailClient ports.EmailClient
	smsClient   ports.SMSClient
}

// NewNotificationService creates a new notification service
func NewNotificationService(emailClient ports.EmailClient, smsClient ports.SMSClient) *NotificationService {
	return &NotificationService{
		emailClient: emailClient,
		smsClient:   smsClient,
	}
}

// SendEmail sends an email notification
func (s *NotificationService) SendEmail(ctx context.Context, to, subject, body string, htmlBody ...string) error {
	if s.emailClient == nil {
		return fmt.Errorf("email client not configured")
	}

	log.Printf("Sending email to %s with subject: %s", to, subject)
	return s.emailClient.SendEmail(ctx, to, subject, body, htmlBody...)
}

// SendSMS sends an SMS notification
func (s *NotificationService) SendSMS(ctx context.Context, phoneNumber, message string) error {
	if s.smsClient == nil {
		return fmt.Errorf("SMS client not configured")
	}

	log.Printf("Sending SMS to %s", phoneNumber)
	return s.smsClient.SendSMS(ctx, phoneNumber, message)
}

// SendNotification sends a notification based on the request type
func (s *NotificationService) SendNotification(ctx context.Context, req *ports.NotificationRequest) error {
	if req == nil {
		return fmt.Errorf("notification request cannot be nil")
	}

	switch req.Type {
	case ports.EmailNotification:
		return s.SendEmail(ctx, req.To, req.Subject, req.Body, req.HTMLBody)
	case ports.SMSNotification:
		return s.SendSMS(ctx, req.PhoneNumber, req.Message)
	default:
		return fmt.Errorf("unsupported notification type: %s", req.Type)
	}
}

// SendBulkEmail sends emails to multiple recipients
func (s *NotificationService) SendBulkEmail(ctx context.Context, recipients []string, subject, body string, htmlBody ...string) error {
	if s.emailClient == nil {
		return fmt.Errorf("email client not configured")
	}

	if len(recipients) == 0 {
		return fmt.Errorf("no recipients provided")
	}

	log.Printf("Sending bulk email to %d recipients with subject: %s", len(recipients), subject)
	return s.emailClient.SendBulkEmail(ctx, recipients, subject, body, htmlBody...)
}

// SendBulkSMS sends SMS to multiple phone numbers
func (s *NotificationService) SendBulkSMS(ctx context.Context, phoneNumbers []string, message string) error {
	if s.smsClient == nil {
		return fmt.Errorf("SMS client not configured")
	}

	if len(phoneNumbers) == 0 {
		return fmt.Errorf("no phone numbers provided")
	}

	log.Printf("Sending bulk SMS to %d recipients", len(phoneNumbers))
	return s.smsClient.SendBulkSMS(ctx, phoneNumbers, message)
}

// ValidateEmail validates if an email address is properly formatted
func (s *NotificationService) ValidateEmail(email string) bool {
	if s.emailClient == nil {
		return false
	}
	return s.emailClient.ValidateEmail(email)
}

// ValidatePhoneNumber validates if a phone number is properly formatted
func (s *NotificationService) ValidatePhoneNumber(phoneNumber string) bool {
	if s.smsClient == nil {
		return false
	}
	return s.smsClient.ValidatePhoneNumber(phoneNumber)
}

// SendOrderConfirmationEmail sends an order confirmation email
func (s *NotificationService) SendOrderConfirmationEmail(ctx context.Context, customerEmail, customerName, orderNumber string, orderItems []OrderItem) error {
	subject := fmt.Sprintf("Order Confirmation - %s", orderNumber)

	// Create plain text body
	body := fmt.Sprintf(`Dear %s,

Thank you for your order! Your order has been confirmed.

Order Number: %s
Order Date: %s

Order Items:
`, customerName, orderNumber, fmt.Sprintf("%d-%02d-%02d", 2024, 1, 1)) // You can get actual date from order

	total := 0.0
	for _, item := range orderItems {
		body += fmt.Sprintf("- %s x%d @ $%.2f = $%.2f\n", item.ProductName, item.Quantity, item.UnitPrice, item.TotalPrice)
		total += item.TotalPrice
	}
	body += fmt.Sprintf("\nTotal: $%.2f\n\nThank you for your business!", total)

	// Create HTML body
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Order Confirmation</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2c3e50;">Order Confirmation</h2>
        
        <p>Dear %s,</p>
        
        <p>Thank you for your order! Your order has been confirmed.</p>
        
        <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <p><strong>Order Number:</strong> %s</p>
            <p><strong>Order Date:</strong> %s</p>
        </div>
        
        <h3>Order Items:</h3>
        <table style="width: 100%%; border-collapse: collapse; margin: 20px 0;">
            <thead>
                <tr style="background-color: #e9ecef;">
                    <th style="padding: 10px; text-align: left; border: 1px solid #dee2e6;">Product</th>
                    <th style="padding: 10px; text-align: center; border: 1px solid #dee2e6;">Quantity</th>
                    <th style="padding: 10px; text-align: right; border: 1px solid #dee2e6;">Price</th>
                    <th style="padding: 10px; text-align: right; border: 1px solid #dee2e6;">Total</th>
                </tr>
            </thead>
            <tbody>`, customerName, orderNumber, fmt.Sprintf("%d-%02d-%02d", 2024, 1, 1))

	for _, item := range orderItems {
		htmlBody += fmt.Sprintf(`
                <tr>
                    <td style="padding: 10px; border: 1px solid #dee2e6;">%s</td>
                    <td style="padding: 10px; text-align: center; border: 1px solid #dee2e6;">%d</td>
                    <td style="padding: 10px; text-align: right; border: 1px solid #dee2e6;">$%.2f</td>
                    <td style="padding: 10px; text-align: right; border: 1px solid #dee2e6;">$%.2f</td>
                </tr>`, item.ProductName, item.Quantity, item.UnitPrice, item.TotalPrice)
	}

	htmlBody += fmt.Sprintf(`
            </tbody>
            <tfoot>
                <tr style="background-color: #f8f9fa; font-weight: bold;">
                    <td colspan="3" style="padding: 10px; text-align: right; border: 1px solid #dee2e6;">Total:</td>
                    <td style="padding: 10px; text-align: right; border: 1px solid #dee2e6;">$%.2f</td>
                </tr>
            </tfoot>
        </table>
        
        <p>Thank you for your business!</p>
        
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="font-size: 12px; color: #666;">This is an automated message. Please do not reply to this email.</p>
    </div>
</body>
</html>`, total)

	return s.SendEmail(ctx, customerEmail, subject, body, htmlBody)
}

// SendOrderConfirmationSMS sends an order confirmation SMS
func (s *NotificationService) SendOrderConfirmationSMS(ctx context.Context, phoneNumber, customerName, orderNumber string, totalAmount float64) error {
	message := fmt.Sprintf("Hi %s! Your order #%s has been confirmed. Total: $%.2f. Thank you for your business!",
		customerName, orderNumber, totalAmount)

	return s.SendSMS(ctx, phoneNumber, message)
}

// SendOrderStatusUpdateSMS sends an order status update SMS
func (s *NotificationService) SendOrderStatusUpdateSMS(ctx context.Context, phoneNumber, customerName, orderNumber, status string) error {
	message := fmt.Sprintf("Hi %s! Your order #%s status has been updated to: %s",
		customerName, orderNumber, status)

	return s.SendSMS(ctx, phoneNumber, message)
}

// OrderItem represents an order item for email templates
type OrderItem struct {
	ProductName string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}
