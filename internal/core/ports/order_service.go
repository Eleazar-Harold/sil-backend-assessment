package ports

import (
	"context"

	"silbackendassessment/internal/core/domain"

	"github.com/google/uuid"
)

// OrderService defines the contract for order business logic
type OrderService interface {
	CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.Order, error)
	GetOrders(ctx context.Context, limit, offset int) ([]*domain.Order, error)
	GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*domain.Order, error)
	GetOrdersByStatus(ctx context.Context, status domain.OrderStatus, limit, offset int) ([]*domain.Order, error)
	UpdateOrder(ctx context.Context, id uuid.UUID, req *domain.UpdateOrderRequest) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
	CancelOrder(ctx context.Context, id uuid.UUID) error
	DeleteOrder(ctx context.Context, id uuid.UUID) error
}
