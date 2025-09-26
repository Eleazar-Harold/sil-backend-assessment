package ports

import (
	"context"

	"silbackendassessment/internal/core/domain"

	"github.com/google/uuid"
)

// CustomerService defines the contract for customer business logic
type CustomerService interface {
	CreateCustomer(ctx context.Context, req *domain.CreateCustomerRequest) (*domain.Customer, error)
	GetCustomer(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	GetCustomers(ctx context.Context, limit, offset int) ([]*domain.Customer, error)
	UpdateCustomer(ctx context.Context, id uuid.UUID, req *domain.UpdateCustomerRequest) (*domain.Customer, error)
	DeleteCustomer(ctx context.Context, id uuid.UUID) error
}
