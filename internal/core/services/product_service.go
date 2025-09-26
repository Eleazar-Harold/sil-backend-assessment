package services

import (
	"context"
	"fmt"
	"time"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
)

type productService struct {
	productRepo  ports.ProductRepository
	categoryRepo ports.CategoryRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo ports.ProductRepository, categoryRepo ports.CategoryRepository) ports.ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error) {
	// Check if product with SKU already exists
	existingProduct, err := s.productRepo.GetBySKU(ctx, req.SKU)
	if err == nil && existingProduct != nil {
		return nil, fmt.Errorf("product with SKU %s already exists", req.SKU)
	}

	// Validate category exists
	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	// Create new product
	product := &domain.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		IsActive:    req.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	return product, nil
}

func (s *productService) GetProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	products, err := s.productRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

func (s *productService) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*domain.Product, error) {
	// Validate category exists
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	products, err := s.productRepo.GetByCategoryID(ctx, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	return products, nil
}

func (s *productService) GetActiveProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	products, err := s.productRepo.GetActiveProducts(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get active products: %w", err)
	}

	return products, nil
}

func (s *productService) SearchProducts(ctx context.Context, name string, limit, offset int) ([]*domain.Product, error) {
	products, err := s.productRepo.SearchByName(ctx, name, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	return products, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req *domain.UpdateProductRequest) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Update fields if provided
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.SKU != nil {
		product.SKU = *req.SKU
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.CategoryID != nil {
		// Validate category exists
		category, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("failed to get category: %w", err)
		}
		if category == nil {
			return nil, fmt.Errorf("category not found")
		}
		product.CategoryID = *req.CategoryID
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *productService) UpdateStock(ctx context.Context, id uuid.UUID, stock int) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return fmt.Errorf("product not found")
	}

	if stock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}

	if err := s.productRepo.UpdateStock(ctx, id, stock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return fmt.Errorf("product not found")
	}

	if err := s.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
