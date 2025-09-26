package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/testutils"

	"github.com/google/uuid"
)

func TestCustomerService_CreateCustomer(t *testing.T) {
	mockRepo := testutils.NewMockCustomerRepository()
	service := NewCustomerService(mockRepo)
	ctx := context.Background()

	t.Run("Create customer successfully", func(t *testing.T) {
		req := &domain.CreateCustomerRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Phone:     "+1234567890",
			Address:   "123 Main St",
			City:      "New York",
			State:     "NY",
			ZipCode:   "10001",
			Country:   "USA",
		}

		customer, err := service.CreateCustomer(ctx, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if customer == nil {
			t.Error("Expected customer to be created")
		}

		if customer.FirstName != "John" {
			t.Errorf("Expected FirstName to be 'John', got: %s", customer.FirstName)
		}

		if customer.LastName != "Doe" {
			t.Errorf("Expected LastName to be 'Doe', got: %s", customer.LastName)
		}

		if customer.Email != "john@example.com" {
			t.Errorf("Expected Email to be 'john@example.com', got: %s", customer.Email)
		}

		if customer.Phone != "+1234567890" {
			t.Errorf("Expected Phone to be '+1234567890', got: %s", customer.Phone)
		}

		if customer.ID == uuid.Nil {
			t.Error("Expected customer ID to be set")
		}

		if customer.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		if customer.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("Create customer with existing email", func(t *testing.T) {
		// Set up existing customer
		existingCustomer := &domain.Customer{
			ID:        uuid.New(),
			FirstName: "Existing",
			LastName:  "Customer",
			Email:     "existing@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.CustomersByEmail["existing@example.com"] = existingCustomer

		req := &domain.CreateCustomerRequest{
			FirstName: "New",
			LastName:  "Customer",
			Email:     "existing@example.com",
		}

		customer, err := service.CreateCustomer(ctx, req)

		if err == nil {
			t.Error("Expected error for existing email")
		}

		if customer != nil {
			t.Error("Expected customer to be nil")
		}

		if err == nil || err.Error() != "customer with email existing@example.com already exists" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("Repository error during creation", func(t *testing.T) {
		mockRepo.CreateError = errors.New("database error")

		req := &domain.CreateCustomerRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		}

		customer, err := service.CreateCustomer(ctx, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if customer != nil {
			t.Error("Expected customer to be nil")
		}
	})
}

func TestCustomerService_GetCustomer(t *testing.T) {
	mockRepo := testutils.NewMockCustomerRepository()
	service := NewCustomerService(mockRepo)
	ctx := context.Background()

	t.Run("Get existing customer", func(t *testing.T) {
		customerID := uuid.New()
		expectedCustomer := &domain.Customer{
			ID:        customerID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Phone:     "+1234567890",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Customers[customerID] = expectedCustomer

		customer, err := service.GetCustomer(ctx, customerID)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if customer == nil {
			t.Error("Expected customer to be returned")
		}

		if customer.ID != customerID {
			t.Errorf("Expected customer ID %s, got: %s", customerID, customer.ID)
		}
	})

	t.Run("Get non-existent customer", func(t *testing.T) {
		customerID := uuid.New()

		customer, err := service.GetCustomer(ctx, customerID)

		if err == nil {
			t.Error("Expected error for non-existent customer")
		}

		if customer != nil {
			t.Error("Expected customer to be nil")
		}

		if !errors.Is(err, testutils.ErrCustomerNotFound) {
			t.Errorf("Expected ErrCustomerNotFound, got: %v", err)
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		customerID := uuid.New()
		mockRepo.GetByIDError = errors.New("database error")

		customer, err := service.GetCustomer(ctx, customerID)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if customer != nil {
			t.Error("Expected customer to be nil")
		}
	})
}

func TestCustomerService_GetCustomers(t *testing.T) {
	mockRepo := testutils.NewMockCustomerRepository()
	service := NewCustomerService(mockRepo)
	ctx := context.Background()

	t.Run("Get customers successfully", func(t *testing.T) {
		customers := []*domain.Customer{
			{
				ID:        uuid.New(),
				FirstName: "Customer 1",
				LastName:  "Last 1",
				Email:     "customer1@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				FirstName: "Customer 2",
				LastName:  "Last 2",
				Email:     "customer2@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRepo.AllCustomers = customers

		result, err := service.GetCustomers(ctx, 10, 0)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 customers, got: %d", len(result))
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo.GetAllError = errors.New("database error")

		customers, err := service.GetCustomers(ctx, 10, 0)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if customers != nil {
			t.Error("Expected customers to be nil")
		}
	})
}

func TestCustomerService_UpdateCustomer(t *testing.T) {
	mockRepo := testutils.NewMockCustomerRepository()
	service := NewCustomerService(mockRepo)
	ctx := context.Background()

	t.Run("Update customer successfully", func(t *testing.T) {
		customerID := uuid.New()
		existingCustomer := &domain.Customer{
			ID:        customerID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Phone:     "+1234567890",
			Address:   "123 Main St",
			City:      "New York",
			State:     "NY",
			ZipCode:   "10001",
			Country:   "USA",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Customers[customerID] = existingCustomer

		newFirstName := "John Updated"
		newPhone := "+0987654321"
		req := &domain.UpdateCustomerRequest{
			FirstName: &newFirstName,
			Phone:     &newPhone,
		}

		customer, err := service.UpdateCustomer(ctx, customerID, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if customer == nil {
			t.Error("Expected customer to be returned")
		}

		if customer.FirstName != "John Updated" {
			t.Errorf("Expected FirstName to be 'John Updated', got: %s", customer.FirstName)
		}

		if customer.Phone != "+0987654321" {
			t.Errorf("Expected Phone to be '+0987654321', got: %s", customer.Phone)
		}

		if customer.Email != "john@example.com" {
			t.Errorf("Expected Email to remain unchanged, got: %s", customer.Email)
		}
	})

	t.Run("Update non-existent customer", func(t *testing.T) {
		customerID := uuid.New()
		newFirstName := "John Updated"
		req := &domain.UpdateCustomerRequest{
			FirstName: &newFirstName,
		}

		customer, err := service.UpdateCustomer(ctx, customerID, req)

		if err == nil {
			t.Error("Expected error for non-existent customer")
		}

		if customer != nil {
			t.Error("Expected customer to be nil")
		}

		if !errors.Is(err, testutils.ErrCustomerNotFound) {
			t.Errorf("Expected ErrCustomerNotFound, got: %v", err)
		}
	})

	t.Run("Update customer email", func(t *testing.T) {
		customerID := uuid.New()
		existingCustomer := &domain.Customer{
			ID:        customerID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Customers[customerID] = existingCustomer

		newEmail := "john.updated@example.com"
		req := &domain.UpdateCustomerRequest{
			Email: &newEmail,
		}

		customer, err := service.UpdateCustomer(ctx, customerID, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if customer.Email != "john.updated@example.com" {
			t.Errorf("Expected Email to be 'john.updated@example.com', got: %s", customer.Email)
		}
	})

	t.Run("Repository error during update", func(t *testing.T) {
		customerID := uuid.New()
		existingCustomer := &domain.Customer{
			ID:        customerID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Customers[customerID] = existingCustomer
		mockRepo.UpdateError = errors.New("database error")

		newFirstName := "John Updated"
		req := &domain.UpdateCustomerRequest{
			FirstName: &newFirstName,
		}

		customer, err := service.UpdateCustomer(ctx, customerID, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if customer != nil {
			t.Error("Expected customer to be nil")
		}
	})
}

func TestCustomerService_DeleteCustomer(t *testing.T) {
	mockRepo := testutils.NewMockCustomerRepository()
	service := NewCustomerService(mockRepo)
	ctx := context.Background()

	t.Run("Delete customer successfully", func(t *testing.T) {
		customerID := uuid.New()
		existingCustomer := &domain.Customer{
			ID:        customerID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Customers[customerID] = existingCustomer

		err := service.DeleteCustomer(ctx, customerID)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify customer was deleted
		if _, exists := mockRepo.Customers[customerID]; exists {
			t.Error("Expected customer to be deleted")
		}
	})

	t.Run("Delete non-existent customer", func(t *testing.T) {
		customerID := uuid.New()

		err := service.DeleteCustomer(ctx, customerID)

		if err == nil {
			t.Error("Expected error for non-existent customer")
		}

		if !errors.Is(err, testutils.ErrCustomerNotFound) {
			t.Errorf("Expected ErrCustomerNotFound, got: %v", err)
		}
	})

	t.Run("Repository error during deletion", func(t *testing.T) {
		customerID := uuid.New()
		existingCustomer := &domain.Customer{
			ID:        customerID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Customers[customerID] = existingCustomer
		mockRepo.DeleteError = errors.New("database error")

		err := service.DeleteCustomer(ctx, customerID)

		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}
