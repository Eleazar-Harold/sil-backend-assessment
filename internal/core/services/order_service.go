package services

import (
	"context"
	"fmt"
	"time"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
)

type orderService struct {
	orderRepo     ports.OrderRepository
	orderItemRepo ports.OrderItemRepository
	customerRepo  ports.CustomerRepository
	productRepo   ports.ProductRepository
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo ports.OrderRepository,
	orderItemRepo ports.OrderItemRepository,
	customerRepo ports.CustomerRepository,
	productRepo ports.ProductRepository,
) ports.OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		customerRepo:  customerRepo,
		productRepo:   productRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.Order, error) {
	// Validate customer exists
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	// Validate products and calculate total
	var totalAmount float64
	var orderItems []*domain.OrderItem

	for _, itemReq := range req.OrderItems {
		// Validate product exists and is active
		product, err := s.productRepo.GetByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product: %w", err)
		}
		if product == nil {
			return nil, fmt.Errorf("product not found: %s", itemReq.ProductID)
		}
		if !product.IsActive {
			return nil, fmt.Errorf("product is not active: %s", product.Name)
		}
		if product.Stock < itemReq.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s: requested %d, available %d", product.Name, itemReq.Quantity, product.Stock)
		}

		// Calculate item total
		itemTotal := product.Price * float64(itemReq.Quantity)
		totalAmount += itemTotal

		// Create order item
		orderItem := &domain.OrderItem{
			ID:         uuid.New(),
			ProductID:  itemReq.ProductID,
			Quantity:   itemReq.Quantity,
			UnitPrice:  product.Price,
			TotalPrice: itemTotal,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		orderItems = append(orderItems, orderItem)
	}

	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%d", time.Now().Unix())

	// Create order
	order := &domain.Order{
		ID:              uuid.New(),
		CustomerID:      req.CustomerID,
		OrderNumber:     orderNumber,
		Status:          domain.OrderStatusPending,
		TotalAmount:     totalAmount,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		Notes:           req.Notes,
		OrderDate:       time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create order in database
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	for _, orderItem := range orderItems {
		orderItem.OrderID = order.ID
		if err := s.orderItemRepo.Create(ctx, orderItem); err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Update product stock
	for _, itemReq := range req.OrderItems {
		product, _ := s.productRepo.GetByID(ctx, itemReq.ProductID)
		newStock := product.Stock - itemReq.Quantity
		if err := s.productRepo.UpdateStock(ctx, itemReq.ProductID, newStock); err != nil {
			return nil, fmt.Errorf("failed to update product stock: %w", err)
		}
	}

	// Set order items for response
	order.OrderItems = make([]domain.OrderItem, len(orderItems))
	for i, item := range orderItems {
		order.OrderItems[i] = *item
	}

	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

func (s *orderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

func (s *orderService) GetOrders(ctx context.Context, limit, offset int) ([]*domain.Order, error) {
	orders, err := s.orderRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return orders, nil
}

func (s *orderService) GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*domain.Order, error) {
	// Validate customer exists
	customer, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	orders, err := s.orderRepo.GetByCustomerID(ctx, customerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by customer: %w", err)
	}

	return orders, nil
}

func (s *orderService) GetOrdersByStatus(ctx context.Context, status domain.OrderStatus, limit, offset int) ([]*domain.Order, error) {
	orders, err := s.orderRepo.GetByStatus(ctx, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by status: %w", err)
	}

	return orders, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, id uuid.UUID, req *domain.UpdateOrderRequest) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	// Update fields if provided
	if req.Status != nil {
		order.Status = *req.Status
	}
	if req.ShippingAddress != nil {
		order.ShippingAddress = *req.ShippingAddress
	}
	if req.BillingAddress != nil {
		order.BillingAddress = *req.BillingAddress
	}
	if req.Notes != nil {
		order.Notes = *req.Notes
	}
	if req.ShippedDate != nil {
		order.ShippedDate = req.ShippedDate
	}
	if req.DeliveredDate != nil {
		order.DeliveredDate = req.DeliveredDate
	}
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order == nil {
		return fmt.Errorf("order not found")
	}

	if err := s.orderRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (s *orderService) CancelOrder(ctx context.Context, id uuid.UUID) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order == nil {
		return fmt.Errorf("order not found")
	}

	// Only allow cancellation of pending or confirmed orders
	if order.Status != domain.OrderStatusPending && order.Status != domain.OrderStatusConfirmed {
		return fmt.Errorf("cannot cancel order with status: %s", order.Status)
	}

	// Restore product stock
	orderItems, err := s.orderItemRepo.GetByOrderID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order items: %w", err)
	}

	for _, item := range orderItems {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product: %w", err)
		}
		newStock := product.Stock + item.Quantity
		if err := s.productRepo.UpdateStock(ctx, item.ProductID, newStock); err != nil {
			return fmt.Errorf("failed to restore product stock: %w", err)
		}
	}

	// Update order status to cancelled
	if err := s.orderRepo.UpdateStatus(ctx, id, domain.OrderStatusCancelled); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

func (s *orderService) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order == nil {
		return fmt.Errorf("order not found")
	}

	// Only allow deletion of cancelled orders
	if order.Status != domain.OrderStatusCancelled {
		return fmt.Errorf("cannot delete order with status: %s", order.Status)
	}

	if err := s.orderRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}
