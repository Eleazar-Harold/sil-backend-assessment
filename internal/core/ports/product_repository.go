package ports

import (
	"context"

	"silbackendassessment/internal/core/domain"

	"github.com/google/uuid"
)

// ProductRepository defines the contract for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetBySKU(ctx context.Context, sku string) (*domain.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	GetByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*domain.Product, error)
	GetActiveProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	SearchByName(ctx context.Context, name string, limit, offset int) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	UpdateStock(ctx context.Context, id uuid.UUID, stock int) error
	Delete(ctx context.Context, id uuid.UUID) error
}
