package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
	"github.com/uptrace/bunrouter"
)

// OrderHandler handles order operations
type OrderHandler struct {
	orderService ports.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService ports.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, req bunrouter.Request) error {
	var createReq domain.CreateOrderRequest
	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	order, err := h.orderService.CreateOrder(req.Context(), &createReq)
	if err != nil {
		http.Error(w, "Failed to create order: "+err.Error(), http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(order)
}

// GetOrder retrieves an order by ID
func (h *OrderHandler) GetOrder(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return err
	}

	order, err := h.orderService.GetOrder(req.Context(), id)
	if err != nil {
		http.Error(w, "Order not found: "+err.Error(), http.StatusNotFound)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(order)
}

// GetOrders retrieves all orders with pagination and filtering
func (h *OrderHandler) GetOrders(w http.ResponseWriter, req bunrouter.Request) error {
	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")
	customerIDStr := req.URL.Query().Get("customer_id")
	statusStr := req.URL.Query().Get("status")

	limit := 10
	offset := 0
	var customerID *uuid.UUID
	var status *domain.OrderStatus

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	if customerIDStr != "" {
		if id, err := uuid.Parse(customerIDStr); err == nil {
			customerID = &id
		}
	}

	if statusStr != "" {
		s := domain.OrderStatus(statusStr)
		status = &s
	}

	var orders []*domain.Order
	var err error

	if customerID != nil {
		orders, err = h.orderService.GetOrdersByCustomer(req.Context(), *customerID, limit, offset)
	} else if status != nil {
		orders, err = h.orderService.GetOrdersByStatus(req.Context(), *status, limit, offset)
	} else {
		orders, err = h.orderService.GetOrders(req.Context(), limit, offset)
	}
	if err != nil {
		http.Error(w, "Failed to get orders: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	response := map[string]interface{}{
		"orders":      orders,
		"limit":       limit,
		"offset":      offset,
		"customer_id": customerID,
		"status":      status,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// UpdateOrder updates an existing order
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return err
	}

	var updateReq domain.UpdateOrderRequest
	if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	order, err := h.orderService.UpdateOrder(req.Context(), id, &updateReq)
	if err != nil {
		http.Error(w, "Failed to update order: "+err.Error(), http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(order)
}

// DeleteOrder deletes an order
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return err
	}

	err = h.orderService.DeleteOrder(req.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete order: "+err.Error(), http.StatusNotFound)
		return err
	}

	response := map[string]string{
		"message": "Order deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registers order routes
func (h *OrderHandler) RegisterRoutes(router *bunrouter.Router) {
	api := router.NewGroup("/api/orders")
	api.POST("", h.CreateOrder)
	api.GET("/:id", h.GetOrder)
	api.GET("", h.GetOrders)
	api.PUT("/:id", h.UpdateOrder)
	api.DELETE("/:id", h.DeleteOrder)
}
