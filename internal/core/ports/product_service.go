package ports

import (
	"context"

	"silbackendassessment/internal/core/domain"

	"github.com/google/uuid"
)

// ProductService defines the contract for product business logic
type ProductService interface {
	CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	GetProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*domain.Product, error)
	GetActiveProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	SearchProducts(ctx context.Context, name string, limit, offset int) ([]*domain.Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req *domain.UpdateProductRequest) (*domain.Product, error)
	UpdateStock(ctx context.Context, id uuid.UUID, stock int) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}
