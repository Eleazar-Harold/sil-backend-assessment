package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ATConfig holds Africa's Talking configuration
type ATConfig struct {
	APIKey   string
	Username string
	BaseURL  string // Default: https://api.africastalking.com
}

// SMSClient handles SMS operations via Africa's Talking
type SMSClient struct {
	config     *ATConfig
	httpClient *http.Client
}

// NewSMSClient creates a new SMS client
func NewSMSClient(config *ATConfig) *SMSClient {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.africastalking.com"
	}

	return &SMSClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendSMS sends an SMS using Africa's Talking API
func (c *SMSClient) SendSMS(ctx context.Context, phoneNumber, message string) error {
	// Validate phone number
	if !c.ValidatePhoneNumber(phoneNumber) {
		return fmt.Errorf("invalid phone number: %s", phoneNumber)
	}

	// Prepare request payload
	payload := SMSRequest{
		Username: c.config.Username,
		Message:  message,
		To:       []string{phoneNumber},
	}

	return c.sendSMSRequest(ctx, payload)
}

// SendBulkSMS sends SMS to multiple phone numbers
func (c *SMSClient) SendBulkSMS(ctx context.Context, phoneNumbers []string, message string) error {
	if len(phoneNumbers) == 0 {
		return fmt.Errorf("no phone numbers provided")
	}

	// Validate all phone numbers
	for _, phoneNumber := range phoneNumbers {
		if !c.ValidatePhoneNumber(phoneNumber) {
			return fmt.Errorf("invalid phone number: %s", phoneNumber)
		}
	}

	// Prepare request payload
	payload := SMSRequest{
		Username: c.config.Username,
		Message:  message,
		To:       phoneNumbers,
	}

	return c.sendSMSRequest(ctx, payload)
}

// ValidatePhoneNumber validates if a phone number is properly formatted
func (c *SMSClient) ValidatePhoneNumber(phoneNumber string) bool {
	// Remove any non-digit characters except +
	cleaned := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phoneNumber, "")

	// Check if it starts with + and has at least 10 digits
	if strings.HasPrefix(cleaned, "+") {
		return len(cleaned) >= 11 && len(cleaned) <= 16
	}

	// Check if it's a local number with at least 10 digits
	return len(cleaned) >= 10 && len(cleaned) <= 15
}

// sendSMSRequest sends the actual SMS request to Africa's Talking API
func (c *SMSClient) sendSMSRequest(ctx context.Context, payload SMSRequest) error {
	// Marshal payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/version1/messaging", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", c.config.APIKey)
	req.Header.Set("Accept", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		var errorResp SMSErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("SMS request failed with status %d and unable to decode error response", resp.StatusCode)
		}
		return fmt.Errorf("SMS request failed: %s", errorResp.ErrorMessage)
	}

	// Parse response
	var smsResp SMSResponse
	if err := json.NewDecoder(resp.Body).Decode(&smsResp); err != nil {
		return fmt.Errorf("failed to decode SMS response: %w", err)
	}

	// Check for any errors in the response
	if len(smsResp.SMSMessageData.Recipients) == 0 {
		return fmt.Errorf("no recipients in SMS response")
	}

	for _, recipient := range smsResp.SMSMessageData.Recipients {
		if recipient.Status != "Success" {
			return fmt.Errorf("SMS delivery failed for %s: %s", recipient.Number, recipient.Status)
		}
	}

	return nil
}

// SMSRequest represents the request payload for Africa's Talking SMS API
type SMSRequest struct {
	Username string   `json:"username"`
	Message  string   `json:"message"`
	To       []string `json:"to"`
}

// SMSResponse represents the response from Africa's Talking SMS API
type SMSResponse struct {
	SMSMessageData struct {
		Message    string `json:"Message"`
		Recipients []struct {
			Number     string `json:"number"`
			Status     string `json:"status"`
			Cost       string `json:"cost"`
			MessageID  string `json:"messageId"`
			StatusCode int    `json:"statusCode"`
			StatusName string `json:"statusName"`
		} `json:"Recipients"`
	} `json:"SMSMessageData"`
}

// SMSErrorResponse represents an error response from Africa's Talking API
type SMSErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}
