package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/testutils"

	"github.com/google/uuid"
)

func TestOrderService_CreateOrder(t *testing.T) {
	mockOrderRepo := testutils.NewMockOrderRepository()
	mockOrderItemRepo := testutils.NewMockOrderItemRepository()
	mockCustomerRepo := testutils.NewMockCustomerRepository()
	mockProductRepo := testutils.NewMockProductRepository()
	service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockCustomerRepo, mockProductRepo)
	ctx := context.Background()

	t.Run("Create order successfully", func(t *testing.T) {
		// Set up customer
		customerID := uuid.New()
		customer := &domain.Customer{
			ID:        customerID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockCustomerRepo.Customers[customerID] = customer

		// Set up products
		productID1 := uuid.New()
		product1 := &domain.Product{
			ID:       productID1,
			Name:     "Product 1",
			SKU:      "PROD-001",
			Price:    99.99,
			Stock:    10,
			IsActive: true,
		}
		mockProductRepo.Products[productID1] = product1

		productID2 := uuid.New()
		product2 := &domain.Product{
			ID:       productID2,
			Name:     "Product 2",
			SKU:      "PROD-002",
			Price:    149.99,
			Stock:    5,
			IsActive: true,
		}
		mockProductRepo.Products[productID2] = product2

		req := &domain.CreateOrderRequest{
			CustomerID: customerID,
			OrderItems: []domain.CreateOrderItemRequest{
				{
					ProductID: productID1,
					Quantity:  2,
				},
				{
					ProductID: productID2,
					Quantity:  1,
				},
			},
			ShippingAddress: "123 Main St, New York, NY 10001",
			BillingAddress:  "123 Main St, New York, NY 10001",
			Notes:           "Test order",
		}

		order, err := service.CreateOrder(ctx, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if order == nil {
			t.Error("Expected order to be created")
		}

		if order.CustomerID != customerID {
			t.Errorf("Expected CustomerID to be %s, got: %s", customerID, order.CustomerID)
		}

		if order.TotalAmount != 349.97 { // (99.99 * 2) + (149.99 * 1)
			t.Errorf("Expected TotalAmount to be 349.97, got: %f", order.TotalAmount)
		}

		if order.Status != "pending" {
			t.Errorf("Expected Status to be 'pending', got: %s", order.Status)
		}

		if order.ShippingAddress != "123 Main St, New York, NY 10001" {
			t.Errorf("Expected ShippingAddress to be correct, got: %s", order.ShippingAddress)
		}

		if order.ID == uuid.Nil {
			t.Error("Expected order ID to be set")
		}

		if order.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		if order.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("Create order with non-existent customer", func(t *testing.T) {
		customerID := uuid.New()

		req := &domain.CreateOrderRequest{
			CustomerID: customerID,
			OrderItems: []domain.CreateOrderItemRequest{
				{
					ProductID: uuid.New(),
					Quantity:  1,
				},
			},
			ShippingAddress: "123 Main St, New York, NY 10001",
			BillingAddress:  "123 Main St, New York, NY 10001",
		}

		order, err := service.CreateOrder(ctx, req)

		if err == nil {
			t.Error("Expected error for non-existent customer")
		}

		if order != nil {
			t.Error("Expected order to be nil")
		}

		if !errors.Is(err, testutils.ErrCustomerNotFound) {
			t.Errorf("Expected ErrCustomerNotFound, got: %v", err)
		}
	})

	t.Run("Create order with non-existent product", func(t *testing.T) {
		// Set up customer
		customerID := uuid.New()
		customer := &domain.Customer{
			ID:        customerID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockCustomerRepo.Customers[customerID] = customer

		req := &domain.CreateOrderRequest{
			CustomerID: customerID,
			OrderItems: []domain.CreateOrderItemRequest{
				{
					ProductID: uuid.New(),
					Quantity:  1,
				},
			},
			ShippingAddress: "123 Main St, New York, NY 10001",
			BillingAddress:  "123 Main St, New York, NY 10001",
		}

		order, err := service.CreateOrder(ctx, req)

		if err == nil {
			t.Error("Expected error for non-existent product")
		}

		if order != nil {
			t.Error("Expected order to be nil")
		}

		if !errors.Is(err, testutils.ErrProductNotFound) {
			t.Errorf("Expected ErrProductNotFound, got: %v", err)
		}
	})

	t.Run("Create order with insufficient stock", func(t *testing.T) {
		// Set up customer
		customerID := uuid.New()
		customer := &domain.Customer{
			ID:        customerID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockCustomerRepo.Customers[customerID] = customer

		// Set up product with low stock
		productID := uuid.New()
		product := &domain.Product{
			ID:       productID,
			Name:     "Product 1",
			SKU:      "PROD-001",
			Price:    99.99,
			Stock:    5, // Only 5 in stock
			IsActive: true,
		}
		mockProductRepo.Products[productID] = product

		req := &domain.CreateOrderRequest{
			CustomerID: customerID,
			OrderItems: []domain.CreateOrderItemRequest{
				{
					ProductID: productID,
					Quantity:  10, // Requesting 10, but only 5 available
				},
			},
			ShippingAddress: "123 Main St, New York, NY 10001",
			BillingAddress:  "123 Main St, New York, NY 10001",
		}

		order, err := service.CreateOrder(ctx, req)

		if err == nil {
			t.Error("Expected error for insufficient stock")
		}

		if order != nil {
			t.Error("Expected order to be nil")
		}

		if err == nil || err.Error() != "insufficient stock for product Product 1: requested 10, available 5" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("Repository error during order creation", func(t *testing.T) {
		// Set up customer
		customerID := uuid.New()
		customer := &domain.Customer{
			ID:        customerID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockCustomerRepo.Customers[customerID] = customer

		// Set up product
		productID := uuid.New()
		product := &domain.Product{
			ID:       productID,
			Name:     "Product 1",
			SKU:      "PROD-001",
			Price:    99.99,
			Stock:    10,
			IsActive: true,
		}
		mockProductRepo.Products[productID] = product
		mockOrderRepo.CreateError = errors.New("database error")

		req := &domain.CreateOrderRequest{
			CustomerID: customerID,
			OrderItems: []domain.CreateOrderItemRequest{
				{
					ProductID: productID,
					Quantity:  1,
				},
			},
			ShippingAddress: "123 Main St, New York, NY 10001",
			BillingAddress:  "123 Main St, New York, NY 10001",
		}

		order, err := service.CreateOrder(ctx, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if order != nil {
			t.Error("Expected order to be nil")
		}
	})
}

func TestOrderService_GetOrder(t *testing.T) {
	mockOrderRepo := testutils.NewMockOrderRepository()
	mockOrderItemRepo := testutils.NewMockOrderItemRepository()
	mockCustomerRepo := testutils.NewMockCustomerRepository()
	mockProductRepo := testutils.NewMockProductRepository()
	service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockCustomerRepo, mockProductRepo)
	ctx := context.Background()

	t.Run("Get existing order", func(t *testing.T) {
		orderID := uuid.New()
		expectedOrder := &domain.Order{
			ID:              orderID,
			CustomerID:      uuid.New(),
			OrderNumber:     "ORD-001",
			Status:          "PENDING",
			TotalAmount:     199.98,
			ShippingAddress: "123 Main St",
			BillingAddress:  "123 Main St",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockOrderRepo.Orders[orderID] = expectedOrder

		order, err := service.GetOrder(ctx, orderID)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if order == nil {
			t.Error("Expected order to be returned")
		}

		if order.ID != orderID {
			t.Errorf("Expected order ID %s, got: %s", orderID, order.ID)
		}
	})

	t.Run("Get non-existent order", func(t *testing.T) {
		orderID := uuid.New()

		order, err := service.GetOrder(ctx, orderID)

		if err == nil {
			t.Error("Expected error for non-existent order")
		}

		if order != nil {
			t.Error("Expected order to be nil")
		}

		if !errors.Is(err, testutils.ErrOrderNotFound) {
			t.Errorf("Expected ErrOrderNotFound, got: %v", err)
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		orderID := uuid.New()
		mockOrderRepo.GetByIDError = errors.New("database error")

		order, err := service.GetOrder(ctx, orderID)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if order != nil {
			t.Error("Expected order to be nil")
		}
	})
}

func TestOrderService_GetOrders(t *testing.T) {
	mockOrderRepo := testutils.NewMockOrderRepository()
	mockOrderItemRepo := testutils.NewMockOrderItemRepository()
	mockCustomerRepo := testutils.NewMockCustomerRepository()
	mockProductRepo := testutils.NewMockProductRepository()
	service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockCustomerRepo, mockProductRepo)
	ctx := context.Background()

	t.Run("Get orders successfully", func(t *testing.T) {
		orders := []*domain.Order{
			{
				ID:          uuid.New(),
				OrderNumber: "ORD-001",
				Status:      "PENDING",
				TotalAmount: 199.98,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          uuid.New(),
				OrderNumber: "ORD-002",
				Status:      "SHIPPED",
				TotalAmount: 299.97,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		mockOrderRepo.AllOrders = orders

		result, err := service.GetOrders(ctx, 10, 0)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 orders, got: %d", len(result))
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		mockOrderRepo.GetAllError = errors.New("database error")

		orders, err := service.GetOrders(ctx, 10, 0)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if orders != nil {
			t.Error("Expected orders to be nil")
		}
	})
}

func TestOrderService_UpdateOrderStatus(t *testing.T) {
	mockOrderRepo := testutils.NewMockOrderRepository()
	mockOrderItemRepo := testutils.NewMockOrderItemRepository()
	mockCustomerRepo := testutils.NewMockCustomerRepository()
	mockProductRepo := testutils.NewMockProductRepository()
	service := NewOrderService(mockOrderRepo, mockOrderItemRepo, mockCustomerRepo, mockProductRepo)
	ctx := context.Background()

	t.Run("Update order status successfully", func(t *testing.T) {
		orderID := uuid.New()
		existingOrder := &domain.Order{
			ID:          orderID,
			OrderNumber: "ORD-001",
			Status:      "PENDING",
			TotalAmount: 199.98,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockOrderRepo.Orders[orderID] = existingOrder

		err := service.UpdateOrderStatus(ctx, orderID, "SHIPPED")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify status was updated
		updatedOrder := mockOrderRepo.Orders[orderID]
		if updatedOrder.Status != "SHIPPED" {
			t.Errorf("Expected status to be 'SHIPPED', got: %s", updatedOrder.Status)
		}
	})

	t.Run("Update non-existent order status", func(t *testing.T) {
		orderID := uuid.New()

		err := service.UpdateOrderStatus(ctx, orderID, "SHIPPED")

		if err == nil {
			t.Error("Expected error for non-existent order")
		}

		if !errors.Is(err, testutils.ErrOrderNotFound) {
			t.Errorf("Expected ErrOrderNotFound, got: %v", err)
		}
	})

	t.Run("Repository error during status update", func(t *testing.T) {
		orderID := uuid.New()
		existingOrder := &domain.Order{
			ID:          orderID,
			OrderNumber: "ORD-001",
			Status:      "PENDING",
			TotalAmount: 199.98,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockOrderRepo.Orders[orderID] = existingOrder
		mockOrderRepo.UpdateError = errors.New("database error")

		err := service.UpdateOrderStatus(ctx, orderID, "SHIPPED")

		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}
