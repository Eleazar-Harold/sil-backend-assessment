package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProduct(t *testing.T) {
	t.Run("Create product with valid data", func(t *testing.T) {
		categoryID := uuid.New()
		product := &Product{
			ID:          uuid.New(),
			Name:        "Test Product",
			Description: "A test product description",
			SKU:         "TEST-SKU-001",
			Price:       99.99,
			Stock:       100,
			CategoryID:  categoryID,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.NotNil(t, product)
		assert.NotEqual(t, uuid.Nil, product.ID)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, "A test product description", product.Description)
		assert.Equal(t, "TEST-SKU-001", product.SKU)
		assert.Equal(t, 99.99, product.Price)
		assert.Equal(t, 100, product.Stock)
		assert.Equal(t, categoryID, product.CategoryID)
		assert.True(t, product.IsActive)
		assert.False(t, product.CreatedAt.IsZero())
		assert.False(t, product.UpdatedAt.IsZero())
	})

	t.Run("Create product with minimal data", func(t *testing.T) {
		categoryID := uuid.New()
		product := &Product{
			ID:         uuid.New(),
			Name:       "Minimal Product",
			SKU:        "MIN-SKU-001",
			Price:      0.0,
			Stock:      0,
			CategoryID: categoryID,
			IsActive:   false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		assert.NotNil(t, product)
		assert.Equal(t, "Minimal Product", product.Name)
		assert.Equal(t, "", product.Description)
		assert.Equal(t, "MIN-SKU-001", product.SKU)
		assert.Equal(t, 0.0, product.Price)
		assert.Equal(t, 0, product.Stock)
		assert.Equal(t, categoryID, product.CategoryID)
		assert.False(t, product.IsActive)
	})

	t.Run("Create product with special characters", func(t *testing.T) {
		categoryID := uuid.New()
		product := &Product{
			ID:          uuid.New(),
			Name:        "Café & Tea Set - 100% Organic",
			Description: "Premium café & tea set with 100% organic ingredients. Perfect for coffee lovers!",
			SKU:         "CAFÉ-TEA-001",
			Price:       149.99,
			Stock:       50,
			CategoryID:  categoryID,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.NotNil(t, product)
		assert.Equal(t, "Café & Tea Set - 100% Organic", product.Name)
		assert.Equal(t, "Premium café & tea set with 100% organic ingredients. Perfect for coffee lovers!", product.Description)
		assert.Equal(t, "CAFÉ-TEA-001", product.SKU)
		assert.Equal(t, 149.99, product.Price)
		assert.Equal(t, 50, product.Stock)
		assert.Equal(t, categoryID, product.CategoryID)
		assert.True(t, product.IsActive)
	})

	t.Run("Create product with high precision price", func(t *testing.T) {
		categoryID := uuid.New()
		product := &Product{
			ID:          uuid.New(),
			Name:        "Precision Product",
			Description: "Product with high precision pricing",
			SKU:         "PREC-001",
			Price:       123.456789,
			Stock:       1,
			CategoryID:  categoryID,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.NotNil(t, product)
		assert.Equal(t, 123.456789, product.Price)
	})
}

func TestCreateProductRequest(t *testing.T) {
	t.Run("Create product request with valid data", func(t *testing.T) {
		categoryID := uuid.New()
		req := &CreateProductRequest{
			Name:        "New Product",
			Description: "A new product description",
			SKU:         "NEW-SKU-001",
			Price:       199.99,
			Stock:       200,
			CategoryID:  categoryID,
			IsActive:    true,
		}

		assert.NotNil(t, req)
		assert.Equal(t, "New Product", req.Name)
		assert.Equal(t, "A new product description", req.Description)
		assert.Equal(t, "NEW-SKU-001", req.SKU)
		assert.Equal(t, 199.99, req.Price)
		assert.Equal(t, 200, req.Stock)
		assert.Equal(t, categoryID, req.CategoryID)
		assert.True(t, req.IsActive)
	})

	t.Run("Create product request with minimal data", func(t *testing.T) {
		categoryID := uuid.New()
		req := &CreateProductRequest{
			Name:       "Minimal Product",
			SKU:        "MIN-SKU-002",
			Price:      0.0,
			Stock:      0,
			CategoryID: categoryID,
			IsActive:   false,
		}

		assert.NotNil(t, req)
		assert.Equal(t, "Minimal Product", req.Name)
		assert.Equal(t, "", req.Description)
		assert.Equal(t, "MIN-SKU-002", req.SKU)
		assert.Equal(t, 0.0, req.Price)
		assert.Equal(t, 0, req.Stock)
		assert.Equal(t, categoryID, req.CategoryID)
		assert.False(t, req.IsActive)
	})

	t.Run("Create product request with negative values", func(t *testing.T) {
		categoryID := uuid.New()
		req := &CreateProductRequest{
			Name:        "Negative Product",
			Description: "Product with negative values",
			SKU:         "NEG-SKU-001",
			Price:       -10.0,
			Stock:       -5,
			CategoryID:  categoryID,
			IsActive:    true,
		}

		assert.NotNil(t, req)
		assert.Equal(t, "Negative Product", req.Name)
		assert.Equal(t, -10.0, req.Price)
		assert.Equal(t, -5, req.Stock)
	})
}

func TestUpdateProductRequest(t *testing.T) {
	t.Run("Update product request with all fields", func(t *testing.T) {
		categoryID := uuid.New()
		newName := "Updated Product"
		newDescription := "Updated description"
		newSKU := "UPD-SKU-001"
		newPrice := 299.99
		newStock := 300
		newIsActive := false

		req := &UpdateProductRequest{
			Name:        &newName,
			Description: &newDescription,
			SKU:         &newSKU,
			Price:       &newPrice,
			Stock:       &newStock,
			CategoryID:  &categoryID,
			IsActive:    &newIsActive,
		}

		assert.NotNil(t, req)
		assert.NotNil(t, req.Name)
		assert.NotNil(t, req.Description)
		assert.NotNil(t, req.SKU)
		assert.NotNil(t, req.Price)
		assert.NotNil(t, req.Stock)
		assert.NotNil(t, req.CategoryID)
		assert.NotNil(t, req.IsActive)
		assert.Equal(t, "Updated Product", *req.Name)
		assert.Equal(t, "Updated description", *req.Description)
		assert.Equal(t, "UPD-SKU-001", *req.SKU)
		assert.Equal(t, 299.99, *req.Price)
		assert.Equal(t, 300, *req.Stock)
		assert.Equal(t, categoryID, *req.CategoryID)
		assert.False(t, *req.IsActive)
	})

	t.Run("Update product request with partial fields", func(t *testing.T) {
		newName := "Partially Updated Product"
		newPrice := 399.99

		req := &UpdateProductRequest{
			Name:  &newName,
			Price: &newPrice,
		}

		assert.NotNil(t, req)
		assert.NotNil(t, req.Name)
		assert.Nil(t, req.Description)
		assert.Nil(t, req.SKU)
		assert.NotNil(t, req.Price)
		assert.Nil(t, req.Stock)
		assert.Nil(t, req.CategoryID)
		assert.Nil(t, req.IsActive)
		assert.Equal(t, "Partially Updated Product", *req.Name)
		assert.Equal(t, 399.99, *req.Price)
	})

	t.Run("Update product request with no fields", func(t *testing.T) {
		req := &UpdateProductRequest{}

		assert.NotNil(t, req)
		assert.Nil(t, req.Name)
		assert.Nil(t, req.Description)
		assert.Nil(t, req.SKU)
		assert.Nil(t, req.Price)
		assert.Nil(t, req.Stock)
		assert.Nil(t, req.CategoryID)
		assert.Nil(t, req.IsActive)
	})

	t.Run("Update product request with zero values", func(t *testing.T) {
		emptyName := ""
		emptyDescription := ""
		emptySKU := ""
		zeroPrice := 0.0
		zeroStock := 0
		zeroCategoryID := uuid.Nil
		falseActive := false

		req := &UpdateProductRequest{
			Name:        &emptyName,
			Description: &emptyDescription,
			SKU:         &emptySKU,
			Price:       &zeroPrice,
			Stock:       &zeroStock,
			CategoryID:  &zeroCategoryID,
			IsActive:    &falseActive,
		}

		assert.NotNil(t, req)
		assert.NotNil(t, req.Name)
		assert.NotNil(t, req.Description)
		assert.NotNil(t, req.SKU)
		assert.NotNil(t, req.Price)
		assert.NotNil(t, req.Stock)
		assert.NotNil(t, req.CategoryID)
		assert.NotNil(t, req.IsActive)
		assert.Equal(t, "", *req.Name)
		assert.Equal(t, "", *req.Description)
		assert.Equal(t, "", *req.SKU)
		assert.Equal(t, 0.0, *req.Price)
		assert.Equal(t, 0, *req.Stock)
		assert.Equal(t, uuid.Nil, *req.CategoryID)
		assert.False(t, *req.IsActive)
	})
}
