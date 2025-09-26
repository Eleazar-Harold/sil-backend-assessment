package ports

import (
	"context"

	"silbackendassessment/internal/core/domain"

	"github.com/google/uuid"
)

// CategoryRepository defines the contract for category data operations
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetByName(ctx context.Context, name string) (*domain.Category, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Category, error)
	GetByParentID(ctx context.Context, parentID uuid.UUID) ([]*domain.Category, error)
	GetRootCategories(ctx context.Context) ([]*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}
