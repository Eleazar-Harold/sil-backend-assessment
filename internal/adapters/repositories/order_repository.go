package repositories

import (
	"context"
	"database/sql"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type orderRepository struct {
	db *bun.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *bun.DB) ports.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	_, err := r.db.NewInsert().Model(order).Exec(ctx)
	return err
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	order := new(domain.Order)
	err := r.db.NewSelect().
		Model(order).
		Relation("Customer").
		Relation("OrderItems").
		Relation("OrderItems.Product").
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	order := new(domain.Order)
	err := r.db.NewSelect().
		Model(order).
		Relation("Customer").
		Relation("OrderItems").
		Relation("OrderItems.Product").
		Where("order_number = ?", orderNumber).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.NewSelect().
		Model(&orders).
		Relation("Customer").
		Where("customer_id = ?", customerID).
		Order("order_date DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return orders, err
}

func (r *orderRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.NewSelect().
		Model(&orders).
		Relation("Customer").
		Order("order_date DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return orders, err
}

func (r *orderRepository) GetByStatus(ctx context.Context, status domain.OrderStatus, limit, offset int) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.NewSelect().
		Model(&orders).
		Relation("Customer").
		Where("status = ?", status).
		Order("order_date DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return orders, err
}

func (r *orderRepository) Update(ctx context.Context, order *domain.Order) error {
	_, err := r.db.NewUpdate().
		Model(order).
		WherePK().
		Exec(ctx)
	return err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Order)(nil)).
		Set("status = ?", status).
		Set("updated_at = CURRENT_TIMESTAMP").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *orderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.Order)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

type orderItemRepository struct {
	db *bun.DB
}

// NewOrderItemRepository creates a new order item repository
func NewOrderItemRepository(db *bun.DB) ports.OrderItemRepository {
	return &orderItemRepository{
		db: db,
	}
}

func (r *orderItemRepository) Create(ctx context.Context, orderItem *domain.OrderItem) error {
	_, err := r.db.NewInsert().Model(orderItem).Exec(ctx)
	return err
}

func (r *orderItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.OrderItem, error) {
	orderItem := new(domain.OrderItem)
	err := r.db.NewSelect().
		Model(orderItem).
		Relation("Order").
		Relation("Product").
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return orderItem, nil
}

func (r *orderItemRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error) {
	var orderItems []*domain.OrderItem
	err := r.db.NewSelect().
		Model(&orderItems).
		Relation("Product").
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Scan(ctx)
	return orderItems, err
}

func (r *orderItemRepository) GetByProductID(ctx context.Context, productID uuid.UUID, limit, offset int) ([]*domain.OrderItem, error) {
	var orderItems []*domain.OrderItem
	err := r.db.NewSelect().
		Model(&orderItems).
		Relation("Order").
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return orderItems, err
}

func (r *orderItemRepository) Update(ctx context.Context, orderItem *domain.OrderItem) error {
	_, err := r.db.NewUpdate().
		Model(orderItem).
		WherePK().
		Exec(ctx)
	return err
}

func (r *orderItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.OrderItem)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *orderItemRepository) DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.OrderItem)(nil)).
		Where("order_id = ?", orderID).
		Exec(ctx)
	return err
}
