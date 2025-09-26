package handlers

import (
	"encoding/json"
	"net/http"

	"silbackendassessment/internal/core/ports"

	"github.com/uptrace/bunrouter"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	notificationService ports.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService ports.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// SendEmailRequest represents a request to send an email
type SendEmailRequest struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	HTMLBody string `json:"html_body,omitempty"`
}

// SendSMSRequest represents a request to send an SMS
type SendSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}

// SendBulkEmailRequest represents a request to send bulk emails
type SendBulkEmailRequest struct {
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	HTMLBody   string   `json:"html_body,omitempty"`
}

// SendBulkSMSRequest represents a request to send bulk SMS
type SendBulkSMSRequest struct {
	PhoneNumbers []string `json:"phone_numbers"`
	Message      string   `json:"message"`
}

// NotificationRequest represents a generic notification request
type NotificationRequest struct {
	Type        string `json:"type"` // "email" or "sms"
	To          string `json:"to,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Body        string `json:"body,omitempty"`
	HTMLBody    string `json:"html_body,omitempty"`
	Message     string `json:"message,omitempty"`
}

// SendEmail sends an email notification
func (h *NotificationHandler) SendEmail(w http.ResponseWriter, req bunrouter.Request) error {
	var request SendEmailRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	// Validate required fields
	if request.To == "" || request.Subject == "" || request.Body == "" {
		http.Error(w, "Missing required fields: to, subject, body", http.StatusBadRequest)
		return nil
	}

	// Send email
	err := h.notificationService.SendEmail(req.Context(), request.To, request.Subject, request.Body, request.HTMLBody)
	if err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{
		"message": "Email sent successfully",
		"to":      request.To,
	})
}

// SendSMS sends an SMS notification
func (h *NotificationHandler) SendSMS(w http.ResponseWriter, req bunrouter.Request) error {
	var request SendSMSRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	// Validate required fields
	if request.PhoneNumber == "" || request.Message == "" {
		http.Error(w, "Missing required fields: phone_number, message", http.StatusBadRequest)
		return nil
	}

	// Send SMS
	err := h.notificationService.SendSMS(req.Context(), request.PhoneNumber, request.Message)
	if err != nil {
		http.Error(w, "Failed to send SMS: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{
		"message":      "SMS sent successfully",
		"phone_number": request.PhoneNumber,
	})
}

// SendBulkEmail sends bulk email notifications
func (h *NotificationHandler) SendBulkEmail(w http.ResponseWriter, req bunrouter.Request) error {
	var request SendBulkEmailRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	// Validate required fields
	if len(request.Recipients) == 0 || request.Subject == "" || request.Body == "" {
		http.Error(w, "Missing required fields: recipients, subject, body", http.StatusBadRequest)
		return nil
	}

	// Send bulk emails
	err := h.notificationService.SendBulkEmail(req.Context(), request.Recipients, request.Subject, request.Body, request.HTMLBody)
	if err != nil {
		http.Error(w, "Failed to send bulk emails: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Bulk emails sent successfully",
		"recipients": request.Recipients,
		"count":      len(request.Recipients),
	})
}

// SendBulkSMS sends bulk SMS notifications
func (h *NotificationHandler) SendBulkSMS(w http.ResponseWriter, req bunrouter.Request) error {
	var request SendBulkSMSRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	// Validate required fields
	if len(request.PhoneNumbers) == 0 || request.Message == "" {
		http.Error(w, "Missing required fields: phone_numbers, message", http.StatusBadRequest)
		return nil
	}

	// Send bulk SMS
	err := h.notificationService.SendBulkSMS(req.Context(), request.PhoneNumbers, request.Message)
	if err != nil {
		http.Error(w, "Failed to send bulk SMS: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Bulk SMS sent successfully",
		"phone_numbers": request.PhoneNumbers,
		"count":         len(request.PhoneNumbers),
	})
}

// SendNotification sends a notification based on type
func (h *NotificationHandler) SendNotification(w http.ResponseWriter, req bunrouter.Request) error {
	var request NotificationRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	// Convert to internal notification request
	notificationReq := &ports.NotificationRequest{
		Type: ports.NotificationType(request.Type),
	}

	switch ports.NotificationType(request.Type) {
	case ports.EmailNotification:
		if request.To == "" || request.Subject == "" || request.Body == "" {
			http.Error(w, "Missing required fields for email: to, subject, body", http.StatusBadRequest)
			return nil
		}
		notificationReq.To = request.To
		notificationReq.Subject = request.Subject
		notificationReq.Body = request.Body
		notificationReq.HTMLBody = request.HTMLBody

	case ports.SMSNotification:
		if request.PhoneNumber == "" || request.Message == "" {
			http.Error(w, "Missing required fields for SMS: phone_number, message", http.StatusBadRequest)
			return nil
		}
		notificationReq.PhoneNumber = request.PhoneNumber
		notificationReq.Message = request.Message

	default:
		http.Error(w, "Invalid notification type. Must be 'email' or 'sms'", http.StatusBadRequest)
		return nil
	}

	// Send notification
	err := h.notificationService.SendNotification(req.Context(), notificationReq)
	if err != nil {
		http.Error(w, "Failed to send notification: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification sent successfully",
		"type":    request.Type,
	})
}

// RegisterRoutes registers all notification routes
func (h *NotificationHandler) RegisterRoutes(router *bunrouter.Router) {
	api := router.NewGroup("/api/notifications")
	api.POST("/email", h.SendEmail)
	api.POST("/sms", h.SendSMS)
	api.POST("/bulk/email", h.SendBulkEmail)
	api.POST("/bulk/sms", h.SendBulkSMS)
	api.POST("/send", h.SendNotification)
}
