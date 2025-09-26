package handlers

import (
	"encoding/json"
	"net/http"

	"silbackendassessment/internal/adapters/middleware"
	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"

	"github.com/uptrace/bunrouter"
)

// CustomerAuthHandler handles customer authentication and profile operations
type CustomerAuthHandler struct {
	customerService ports.CustomerService
	authService     ports.AuthService
}

// NewCustomerAuthHandler creates a new customer auth handler
func NewCustomerAuthHandler(
	customerService ports.CustomerService,
	authService ports.AuthService,
) *CustomerAuthHandler {
	return &CustomerAuthHandler{
		customerService: customerService,
		authService:     authService,
	}
}

// ProfileResponse represents the customer profile response
type ProfileResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UpdateProfileRequest represents the request to update customer profile
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Address   *string `json:"address,omitempty"`
	City      *string `json:"city,omitempty"`
	State     *string `json:"state,omitempty"`
	ZipCode   *string `json:"zip_code,omitempty"`
	Country   *string `json:"country,omitempty"`
}

// GetProfile returns the current customer's profile
func (h *CustomerAuthHandler) GetProfile(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	// Get customer from context (set by auth middleware)
	customerInfo, ok := middleware.GetCustomerFromContext(ctx)
	if !ok {
		http.Error(w, "Customer not authenticated", http.StatusUnauthorized)
		return nil
	}

	// Get customer details from service
	customerID, err := uuid.Parse(customerInfo.ID)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return err
	}
	customer, err := h.customerService.GetCustomer(ctx, customerID)
	if err != nil {
		http.Error(w, "Failed to get customer profile: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	response := ProfileResponse{
		ID:        customer.ID.String(),
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Address:   customer.Address,
		City:      customer.City,
		State:     customer.State,
		ZipCode:   customer.ZipCode,
		Country:   customer.Country,
		CreatedAt: customer.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: customer.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// UpdateProfile updates the current customer's profile
func (h *CustomerAuthHandler) UpdateProfile(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	// Get customer from context (set by auth middleware)
	customerInfo, ok := middleware.GetCustomerFromContext(ctx)
	if !ok {
		http.Error(w, "Customer not authenticated", http.StatusUnauthorized)
		return nil
	}

	// Parse request body
	var updateReq UpdateProfileRequest
	if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	// Convert to domain request
	domainReq := &domain.UpdateCustomerRequest{
		FirstName: updateReq.FirstName,
		LastName:  updateReq.LastName,
		Phone:     updateReq.Phone,
		Address:   updateReq.Address,
		City:      updateReq.City,
		State:     updateReq.State,
		ZipCode:   updateReq.ZipCode,
		Country:   updateReq.Country,
	}

	// Update customer
	customerID, err := uuid.Parse(customerInfo.ID)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return err
	}
	customer, err := h.customerService.UpdateCustomer(ctx, customerID, domainReq)
	if err != nil {
		http.Error(w, "Failed to update customer profile: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	response := ProfileResponse{
		ID:        customer.ID.String(),
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Address:   customer.Address,
		City:      customer.City,
		State:     customer.State,
		ZipCode:   customer.ZipCode,
		Country:   customer.Country,
		CreatedAt: customer.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: customer.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// DeleteAccount deletes the current customer's account
func (h *CustomerAuthHandler) DeleteAccount(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	// Get customer from context (set by auth middleware)
	customerInfo, ok := middleware.GetCustomerFromContext(ctx)
	if !ok {
		http.Error(w, "Customer not authenticated", http.StatusUnauthorized)
		return nil
	}

	// Delete customer
	customerID, err := uuid.Parse(customerInfo.ID)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return err
	}
	err = h.customerService.DeleteCustomer(ctx, customerID)
	if err != nil {
		http.Error(w, "Failed to delete customer account: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	response := map[string]string{
		"message": "Account deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registers customer auth routes
func (h *CustomerAuthHandler) RegisterRoutes(router *bunrouter.Router, authMiddleware *middleware.AuthMiddleware) {
	// Protected routes that require customer authentication
	protected := router.NewGroup("/api/customer")
	protected.Use(authMiddleware.RequireCustomerAuth)

	protected.GET("/profile", h.GetProfile)
	protected.PUT("/profile", h.UpdateProfile)
	protected.DELETE("/account", h.DeleteAccount)
}
