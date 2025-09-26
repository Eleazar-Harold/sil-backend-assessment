package ports

import (
	"context"

	"silbackendassessment/internal/core/domain"

	"github.com/google/uuid"
)

// OrderRepository defines the contract for order data operations
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*domain.Order, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Order, error)
	GetByStatus(ctx context.Context, status domain.OrderStatus, limit, offset int) ([]*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// OrderItemRepository defines the contract for order item data operations
type OrderItemRepository interface {
	Create(ctx context.Context, orderItem *domain.OrderItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.OrderItem, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error)
	GetByProductID(ctx context.Context, productID uuid.UUID, limit, offset int) ([]*domain.OrderItem, error)
	Update(ctx context.Context, orderItem *domain.OrderItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error
}
