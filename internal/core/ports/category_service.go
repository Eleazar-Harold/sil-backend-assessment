package ports

import (
	"context"

	"silbackendassessment/internal/core/domain"

	"github.com/google/uuid"
)

// CategoryService defines the contract for category business logic
type CategoryService interface {
	CreateCategory(ctx context.Context, req *domain.CreateCategoryRequest) (*domain.Category, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetCategories(ctx context.Context, limit, offset int) ([]*domain.Category, error)
	GetRootCategories(ctx context.Context) ([]*domain.Category, error)
	GetSubCategories(ctx context.Context, parentID uuid.UUID) ([]*domain.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req *domain.UpdateCategoryRequest) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}
