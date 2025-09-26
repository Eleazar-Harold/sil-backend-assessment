package ports

import (
	"context"

	"silbackendassessment/internal/core/domain"

	"github.com/google/uuid"
)

// CustomerRepository defines the contract for customer data operations
type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	GetByEmail(ctx context.Context, email string) (*domain.Customer, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Customer, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}
