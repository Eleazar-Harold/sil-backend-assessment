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

func TestProductService_CreateProduct(t *testing.T) {
	mockProductRepo := testutils.NewMockProductRepository()
	mockCategoryRepo := testutils.NewMockCategoryRepository()
	service := NewProductService(mockProductRepo, mockCategoryRepo)
	ctx := context.Background()

	t.Run("Create product successfully", func(t *testing.T) {
		// Set up category
		categoryID := uuid.New()
		category := &domain.Category{
			ID:          categoryID,
			Name:        "Electronics",
			Description: "Electronic products",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockCategoryRepo.Categories[categoryID] = category

		req := &domain.CreateProductRequest{
			Name:        "Test Product",
			Description: "A test product",
			SKU:         "TEST-001",
			Price:       99.99,
			Stock:       10,
			CategoryID:  categoryID,
		}

		product, err := service.CreateProduct(ctx, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if product == nil {
			t.Error("Expected product to be created")
		}

		if product.Name != "Test Product" {
			t.Errorf("Expected Name to be 'Test Product', got: %s", product.Name)
		}

		if product.SKU != "TEST-001" {
			t.Errorf("Expected SKU to be 'TEST-001', got: %s", product.SKU)
		}

		if product.Price != 99.99 {
			t.Errorf("Expected Price to be 99.99, got: %f", product.Price)
		}

		if product.Stock != 10 {
			t.Errorf("Expected Stock to be 10, got: %d", product.Stock)
		}

		if product.CategoryID != categoryID {
			t.Errorf("Expected CategoryID to be %s, got: %s", categoryID, product.CategoryID)
		}

		if product.ID == uuid.Nil {
			t.Error("Expected product ID to be set")
		}

		if product.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		if product.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("Create product with existing SKU", func(t *testing.T) {
		// Set up existing product
		existingProduct := &domain.Product{
			ID:        uuid.New(),
			Name:      "Existing Product",
			SKU:       "EXISTING-001",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockProductRepo.ProductsBySKU["EXISTING-001"] = existingProduct

		// Set up category
		categoryID := uuid.New()
		category := &domain.Category{
			ID:          categoryID,
			Name:        "Electronics",
			Description: "Electronic products",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockCategoryRepo.Categories[categoryID] = category

		req := &domain.CreateProductRequest{
			Name:        "New Product",
			Description: "A new product",
			SKU:         "EXISTING-001",
			Price:       99.99,
			Stock:       10,
			CategoryID:  categoryID,
		}

		product, err := service.CreateProduct(ctx, req)

		if err == nil {
			t.Error("Expected error for existing SKU")
		}

		if product != nil {
			t.Error("Expected product to be nil")
		}

		if err == nil || err.Error() != "product with SKU EXISTING-001 already exists" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("Create product with non-existent category", func(t *testing.T) {
		categoryID := uuid.New()

		req := &domain.CreateProductRequest{
			Name:        "Test Product",
			Description: "A test product",
			SKU:         "TEST-002",
			Price:       99.99,
			Stock:       10,
			CategoryID:  categoryID,
		}

		product, err := service.CreateProduct(ctx, req)

		if err == nil {
			t.Error("Expected error for non-existent category")
		}

		if product != nil {
			t.Error("Expected product to be nil")
		}

		if !errors.Is(err, testutils.ErrCategoryNotFound) {
			t.Errorf("Expected ErrCategoryNotFound, got: %v", err)
		}
	})

	t.Run("Repository error during creation", func(t *testing.T) {
		// Set up category
		categoryID := uuid.New()
		category := &domain.Category{
			ID:          categoryID,
			Name:        "Electronics",
			Description: "Electronic products",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockCategoryRepo.Categories[categoryID] = category
		mockProductRepo.CreateError = errors.New("database error")

		req := &domain.CreateProductRequest{
			Name:        "Test Product",
			Description: "A test product",
			SKU:         "TEST-003",
			Price:       99.99,
			Stock:       10,
			CategoryID:  categoryID,
		}

		product, err := service.CreateProduct(ctx, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if product != nil {
			t.Error("Expected product to be nil")
		}
	})
}

func TestProductService_GetProduct(t *testing.T) {
	mockProductRepo := testutils.NewMockProductRepository()
	mockCategoryRepo := testutils.NewMockCategoryRepository()
	service := NewProductService(mockProductRepo, mockCategoryRepo)
	ctx := context.Background()

	t.Run("Get existing product", func(t *testing.T) {
		productID := uuid.New()
		expectedProduct := &domain.Product{
			ID:        productID,
			Name:      "Test Product",
			SKU:       "TEST-001",
			Price:     99.99,
			Stock:     10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockProductRepo.Products[productID] = expectedProduct

		product, err := service.GetProduct(ctx, productID)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if product == nil {
			t.Error("Expected product to be returned")
		}

		if product.ID != productID {
			t.Errorf("Expected product ID %s, got: %s", productID, product.ID)
		}
	})

	t.Run("Get non-existent product", func(t *testing.T) {
		productID := uuid.New()

		product, err := service.GetProduct(ctx, productID)

		if err == nil {
			t.Error("Expected error for non-existent product")
		}

		if product != nil {
			t.Error("Expected product to be nil")
		}

		if !errors.Is(err, testutils.ErrProductNotFound) {
			t.Errorf("Expected ErrProductNotFound, got: %v", err)
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		productID := uuid.New()
		mockProductRepo.GetByIDError = errors.New("database error")

		product, err := service.GetProduct(ctx, productID)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if product != nil {
			t.Error("Expected product to be nil")
		}
	})
}

func TestProductService_GetProducts(t *testing.T) {
	mockProductRepo := testutils.NewMockProductRepository()
	mockCategoryRepo := testutils.NewMockCategoryRepository()
	service := NewProductService(mockProductRepo, mockCategoryRepo)
	ctx := context.Background()

	t.Run("Get products successfully", func(t *testing.T) {
		products := []*domain.Product{
			{
				ID:        uuid.New(),
				Name:      "Product 1",
				SKU:       "PROD-001",
				Price:     99.99,
				Stock:     10,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				Name:      "Product 2",
				SKU:       "PROD-002",
				Price:     199.99,
				Stock:     5,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockProductRepo.AllProducts = products

		result, err := service.GetProducts(ctx, 10, 0)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 products, got: %d", len(result))
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		mockProductRepo.GetAllError = errors.New("database error")

		products, err := service.GetProducts(ctx, 10, 0)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if products != nil {
			t.Error("Expected products to be nil")
		}
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	mockProductRepo := testutils.NewMockProductRepository()
	mockCategoryRepo := testutils.NewMockCategoryRepository()
	service := NewProductService(mockProductRepo, mockCategoryRepo)
	ctx := context.Background()

	t.Run("Update product successfully", func(t *testing.T) {
		productID := uuid.New()
		existingProduct := &domain.Product{
			ID:        productID,
			Name:      "Original Product",
			SKU:       "ORIG-001",
			Price:     99.99,
			Stock:     10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockProductRepo.Products[productID] = existingProduct

		newName := "Updated Product"
		newPrice := 149.99
		req := &domain.UpdateProductRequest{
			Name:  &newName,
			Price: &newPrice,
		}

		product, err := service.UpdateProduct(ctx, productID, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if product == nil {
			t.Error("Expected product to be returned")
		}

		if product.Name != "Updated Product" {
			t.Errorf("Expected Name to be 'Updated Product', got: %s", product.Name)
		}

		if product.Price != 149.99 {
			t.Errorf("Expected Price to be 149.99, got: %f", product.Price)
		}

		if product.SKU != "ORIG-001" {
			t.Errorf("Expected SKU to remain unchanged, got: %s", product.SKU)
		}
	})

	t.Run("Update non-existent product", func(t *testing.T) {
		productID := uuid.New()
		newName := "Updated Product"
		req := &domain.UpdateProductRequest{
			Name: &newName,
		}

		product, err := service.UpdateProduct(ctx, productID, req)

		if err == nil {
			t.Error("Expected error for non-existent product")
		}

		if product != nil {
			t.Error("Expected product to be nil")
		}

		if !errors.Is(err, testutils.ErrProductNotFound) {
			t.Errorf("Expected ErrProductNotFound, got: %v", err)
		}
	})

	t.Run("Repository error during update", func(t *testing.T) {
		productID := uuid.New()
		existingProduct := &domain.Product{
			ID:        productID,
			Name:      "Original Product",
			SKU:       "ORIG-001",
			Price:     99.99,
			Stock:     10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockProductRepo.Products[productID] = existingProduct
		mockProductRepo.UpdateError = errors.New("database error")

		newName := "Updated Product"
		req := &domain.UpdateProductRequest{
			Name: &newName,
		}

		product, err := service.UpdateProduct(ctx, productID, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if product != nil {
			t.Error("Expected product to be nil")
		}
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	mockProductRepo := testutils.NewMockProductRepository()
	mockCategoryRepo := testutils.NewMockCategoryRepository()
	service := NewProductService(mockProductRepo, mockCategoryRepo)
	ctx := context.Background()

	t.Run("Delete product successfully", func(t *testing.T) {
		productID := uuid.New()
		existingProduct := &domain.Product{
			ID:        productID,
			Name:      "Test Product",
			SKU:       "TEST-001",
			Price:     99.99,
			Stock:     10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockProductRepo.Products[productID] = existingProduct

		err := service.DeleteProduct(ctx, productID)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify product was deleted
		if _, exists := mockProductRepo.Products[productID]; exists {
			t.Error("Expected product to be deleted")
		}
	})

	t.Run("Delete non-existent product", func(t *testing.T) {
		productID := uuid.New()

		err := service.DeleteProduct(ctx, productID)

		if err == nil {
			t.Error("Expected error for non-existent product")
		}

		if !errors.Is(err, testutils.ErrProductNotFound) {
			t.Errorf("Expected ErrProductNotFound, got: %v", err)
		}
	})

	t.Run("Repository error during deletion", func(t *testing.T) {
		productID := uuid.New()
		existingProduct := &domain.Product{
			ID:        productID,
			Name:      "Test Product",
			SKU:       "TEST-001",
			Price:     99.99,
			Stock:     10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockProductRepo.Products[productID] = existingProduct
		mockProductRepo.DeleteError = errors.New("database error")

		err := service.DeleteProduct(ctx, productID)

		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}
