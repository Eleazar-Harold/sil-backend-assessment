package notifications

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	TLS      bool
}

// EmailClient handles email operations via SMTP
type EmailClient struct {
	config *SMTPConfig
}

// NewEmailClient creates a new email client
func NewEmailClient(config *SMTPConfig) *EmailClient {
	return &EmailClient{
		config: config,
	}
}

// SendEmail sends an email using SMTP
func (c *EmailClient) SendEmail(ctx context.Context, to, subject, body string, htmlBody ...string) error {
	// Validate email addresses
	if !c.ValidateEmail(to) {
		return fmt.Errorf("invalid recipient email address: %s", to)
	}
	if !c.ValidateEmail(c.config.From) {
		return fmt.Errorf("invalid sender email address: %s", c.config.From)
	}

	// Prepare message
	message := c.buildMessage(to, subject, body, htmlBody...)

	// Set up authentication
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)

	// Build server address
	serverAddr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	// Send email
	if c.config.TLS {
		return c.sendWithTLS(ctx, serverAddr, auth, c.config.From, []string{to}, []byte(message))
	}

	return smtp.SendMail(serverAddr, auth, c.config.From, []string{to}, []byte(message))
}

// SendBulkEmail sends emails to multiple recipients
func (c *EmailClient) SendBulkEmail(ctx context.Context, recipients []string, subject, body string, htmlBody ...string) error {
	if len(recipients) == 0 {
		return fmt.Errorf("no recipients provided")
	}

	// Validate all email addresses
	for _, recipient := range recipients {
		if !c.ValidateEmail(recipient) {
			return fmt.Errorf("invalid email address: %s", recipient)
		}
	}

	// Prepare message
	message := c.buildMessage("", subject, body, htmlBody...)

	// Set up authentication
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)

	// Build server address
	serverAddr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	// Send emails
	if c.config.TLS {
		return c.sendBulkWithTLS(ctx, serverAddr, auth, c.config.From, recipients, []byte(message))
	}

	return smtp.SendMail(serverAddr, auth, c.config.From, recipients, []byte(message))
}

// ValidateEmail validates if an email address is properly formatted
func (c *EmailClient) ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// buildMessage builds the email message with headers and body
func (c *EmailClient) buildMessage(to, subject, body string, htmlBody ...string) string {
	headers := make(map[string]string)
	headers["From"] = c.config.From
	if to != "" {
		headers["To"] = to
	}
	headers["Subject"] = subject
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	headers["MIME-Version"] = "1.0"

	var message strings.Builder

	// Write headers
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Handle HTML content
	if len(htmlBody) > 0 && htmlBody[0] != "" {
		boundary := fmt.Sprintf("boundary_%d", time.Now().Unix())
		headers["Content-Type"] = fmt.Sprintf("multipart/alternative; boundary=\"%s\"", boundary)

		message.Reset()
		for key, value := range headers {
			message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
		message.WriteString("\r\n")

		// Plain text part
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
		message.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
		message.WriteString(body)
		message.WriteString("\r\n")

		// HTML part
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		message.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
		message.WriteString(htmlBody[0])
		message.WriteString("\r\n")

		message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		headers["Content-Type"] = "text/plain; charset=\"UTF-8\""
		headers["Content-Transfer-Encoding"] = "8bit"

		message.Reset()
		for key, value := range headers {
			message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
		message.WriteString("\r\n")
		message.WriteString(body)
	}

	return message.String()
}

// sendWithTLS sends email with TLS encryption
func (c *EmailClient) sendWithTLS(ctx context.Context, serverAddr string, auth smtp.Auth, from string, to []string, message []byte) error {
	// Create TLS connection
	conn, err := tls.Dial("tcp", serverAddr, &tls.Config{ServerName: c.config.Host})
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, c.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Set sender
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	if _, err := writer.Write(message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return nil
}

// sendBulkWithTLS sends bulk emails with TLS encryption
func (c *EmailClient) sendBulkWithTLS(ctx context.Context, serverAddr string, auth smtp.Auth, from string, recipients []string, message []byte) error {
	// Create TLS connection
	conn, err := tls.Dial("tcp", serverAddr, &tls.Config{ServerName: c.config.Host})
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, c.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Set sender
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	if _, err := writer.Write(message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return nil
}
