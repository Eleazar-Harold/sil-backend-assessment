package ports

import (
	"context"

	"github.com/google/uuid"
	"silbackendassessment/internal/core/domain"
)

// UserService defines the contract for user business logic
type UserService interface {
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}