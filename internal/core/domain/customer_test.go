package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCustomer(t *testing.T) {
	t.Run("Create customer with valid data", func(t *testing.T) {
		customer := &Customer{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone:     "+1234567890",
			Address:   "123 Main St",
			City:      "New York",
			State:     "NY",
			ZipCode:   "10001",
			Country:   "USA",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, customer)
		assert.NotEqual(t, uuid.Nil, customer.ID)
		assert.Equal(t, "John", customer.FirstName)
		assert.Equal(t, "Doe", customer.LastName)
		assert.Equal(t, "john.doe@example.com", customer.Email)
		assert.Equal(t, "+1234567890", customer.Phone)
		assert.Equal(t, "123 Main St", customer.Address)
		assert.Equal(t, "New York", customer.City)
		assert.Equal(t, "NY", customer.State)
		assert.Equal(t, "10001", customer.ZipCode)
		assert.Equal(t, "USA", customer.Country)
		assert.False(t, customer.CreatedAt.IsZero())
		assert.False(t, customer.UpdatedAt.IsZero())
	})

	t.Run("Create customer with minimal data", func(t *testing.T) {
		customer := &Customer{
			ID:        uuid.New(),
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, customer)
		assert.Equal(t, "Jane", customer.FirstName)
		assert.Equal(t, "Smith", customer.LastName)
		assert.Equal(t, "jane.smith@example.com", customer.Email)
		assert.Equal(t, "", customer.Phone)
		assert.Equal(t, "", customer.Address)
		assert.Equal(t, "", customer.City)
		assert.Equal(t, "", customer.State)
		assert.Equal(t, "", customer.ZipCode)
		assert.Equal(t, "", customer.Country)
	})

	t.Run("Create customer with special characters", func(t *testing.T) {
		customer := &Customer{
			ID:        uuid.New(),
			FirstName: "José",
			LastName:  "García-López",
			Email:     "josé.garcía+test@example.com",
			Phone:     "+34 123 456 789",
			Address:   "Calle Mayor, 123",
			City:      "Madrid",
			State:     "Madrid",
			ZipCode:   "28001",
			Country:   "España",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, customer)
		assert.Equal(t, "José", customer.FirstName)
		assert.Equal(t, "García-López", customer.LastName)
		assert.Equal(t, "josé.garcía+test@example.com", customer.Email)
		assert.Equal(t, "+34 123 456 789", customer.Phone)
		assert.Equal(t, "Calle Mayor, 123", customer.Address)
		assert.Equal(t, "Madrid", customer.City)
		assert.Equal(t, "Madrid", customer.State)
		assert.Equal(t, "28001", customer.ZipCode)
		assert.Equal(t, "España", customer.Country)
	})
}

func TestCreateCustomerRequest(t *testing.T) {
	t.Run("Create customer request with valid data", func(t *testing.T) {
		req := &CreateCustomerRequest{
			FirstName: "Alice",
			LastName:  "Johnson",
			Email:     "alice.johnson@example.com",
			Phone:     "+1987654321",
			Address:   "456 Oak Ave",
			City:      "Los Angeles",
			State:     "CA",
			ZipCode:   "90210",
			Country:   "USA",
		}

		assert.NotNil(t, req)
		assert.Equal(t, "Alice", req.FirstName)
		assert.Equal(t, "Johnson", req.LastName)
		assert.Equal(t, "alice.johnson@example.com", req.Email)
		assert.Equal(t, "+1987654321", req.Phone)
		assert.Equal(t, "456 Oak Ave", req.Address)
		assert.Equal(t, "Los Angeles", req.City)
		assert.Equal(t, "CA", req.State)
		assert.Equal(t, "90210", req.ZipCode)
		assert.Equal(t, "USA", req.Country)
	})

	t.Run("Create customer request with minimal data", func(t *testing.T) {
		req := &CreateCustomerRequest{
			FirstName: "Bob",
			LastName:  "Wilson",
			Email:     "bob.wilson@example.com",
		}

		assert.NotNil(t, req)
		assert.Equal(t, "Bob", req.FirstName)
		assert.Equal(t, "Wilson", req.LastName)
		assert.Equal(t, "bob.wilson@example.com", req.Email)
		assert.Equal(t, "", req.Phone)
		assert.Equal(t, "", req.Address)
		assert.Equal(t, "", req.City)
		assert.Equal(t, "", req.State)
		assert.Equal(t, "", req.ZipCode)
		assert.Equal(t, "", req.Country)
	})

	t.Run("Create customer request with empty required fields", func(t *testing.T) {
		req := &CreateCustomerRequest{
			FirstName: "",
			LastName:  "",
			Email:     "",
		}

		assert.NotNil(t, req)
		assert.Equal(t, "", req.FirstName)
		assert.Equal(t, "", req.LastName)
		assert.Equal(t, "", req.Email)
	})
}

func TestUpdateCustomerRequest(t *testing.T) {
	t.Run("Update customer request with all fields", func(t *testing.T) {
		newFirstName := "Updated First"
		newLastName := "Updated Last"
		newEmail := "updated@example.com"
		newPhone := "+1111111111"
		newAddress := "Updated Address"
		newCity := "Updated City"
		newState := "Updated State"
		newZipCode := "Updated Zip"
		newCountry := "Updated Country"

		req := &UpdateCustomerRequest{
			FirstName: &newFirstName,
			LastName:  &newLastName,
			Email:     &newEmail,
			Phone:     &newPhone,
			Address:   &newAddress,
			City:      &newCity,
			State:     &newState,
			ZipCode:   &newZipCode,
			Country:   &newCountry,
		}

		assert.NotNil(t, req)
		assert.NotNil(t, req.FirstName)
		assert.NotNil(t, req.LastName)
		assert.NotNil(t, req.Email)
		assert.NotNil(t, req.Phone)
		assert.NotNil(t, req.Address)
		assert.NotNil(t, req.City)
		assert.NotNil(t, req.State)
		assert.NotNil(t, req.ZipCode)
		assert.NotNil(t, req.Country)
		assert.Equal(t, "Updated First", *req.FirstName)
		assert.Equal(t, "Updated Last", *req.LastName)
		assert.Equal(t, "updated@example.com", *req.Email)
		assert.Equal(t, "+1111111111", *req.Phone)
		assert.Equal(t, "Updated Address", *req.Address)
		assert.Equal(t, "Updated City", *req.City)
		assert.Equal(t, "Updated State", *req.State)
		assert.Equal(t, "Updated Zip", *req.ZipCode)
		assert.Equal(t, "Updated Country", *req.Country)
	})

	t.Run("Update customer request with partial fields", func(t *testing.T) {
		newFirstName := "Partially Updated"
		newEmail := "partial@example.com"

		req := &UpdateCustomerRequest{
			FirstName: &newFirstName,
			Email:     &newEmail,
		}

		assert.NotNil(t, req)
		assert.NotNil(t, req.FirstName)
		assert.Nil(t, req.LastName)
		assert.NotNil(t, req.Email)
		assert.Nil(t, req.Phone)
		assert.Nil(t, req.Address)
		assert.Nil(t, req.City)
		assert.Nil(t, req.State)
		assert.Nil(t, req.ZipCode)
		assert.Nil(t, req.Country)
		assert.Equal(t, "Partially Updated", *req.FirstName)
		assert.Equal(t, "partial@example.com", *req.Email)
	})

	t.Run("Update customer request with no fields", func(t *testing.T) {
		req := &UpdateCustomerRequest{}

		assert.NotNil(t, req)
		assert.Nil(t, req.FirstName)
		assert.Nil(t, req.LastName)
		assert.Nil(t, req.Email)
		assert.Nil(t, req.Phone)
		assert.Nil(t, req.Address)
		assert.Nil(t, req.City)
		assert.Nil(t, req.State)
		assert.Nil(t, req.ZipCode)
		assert.Nil(t, req.Country)
	})

	t.Run("Update customer request with empty string values", func(t *testing.T) {
		emptyFirstName := ""
		emptyLastName := ""
		emptyEmail := ""
		emptyPhone := ""
		emptyAddress := ""
		emptyCity := ""
		emptyState := ""
		emptyZipCode := ""
		emptyCountry := ""

		req := &UpdateCustomerRequest{
			FirstName: &emptyFirstName,
			LastName:  &emptyLastName,
			Email:     &emptyEmail,
			Phone:     &emptyPhone,
			Address:   &emptyAddress,
			City:      &emptyCity,
			State:     &emptyState,
			ZipCode:   &emptyZipCode,
			Country:   &emptyCountry,
		}

		assert.NotNil(t, req)
		assert.NotNil(t, req.FirstName)
		assert.NotNil(t, req.LastName)
		assert.NotNil(t, req.Email)
		assert.NotNil(t, req.Phone)
		assert.NotNil(t, req.Address)
		assert.NotNil(t, req.City)
		assert.NotNil(t, req.State)
		assert.NotNil(t, req.ZipCode)
		assert.NotNil(t, req.Country)
		assert.Equal(t, "", *req.FirstName)
		assert.Equal(t, "", *req.LastName)
		assert.Equal(t, "", *req.Email)
		assert.Equal(t, "", *req.Phone)
		assert.Equal(t, "", *req.Address)
		assert.Equal(t, "", *req.City)
		assert.Equal(t, "", *req.State)
		assert.Equal(t, "", *req.ZipCode)
		assert.Equal(t, "", *req.Country)
	})
}
