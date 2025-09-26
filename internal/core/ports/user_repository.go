package ports

import (
	"context"

	"github.com/google/uuid"
	"silbackendassessment/internal/core/domain"
)

// UserRepository defines the contract for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}