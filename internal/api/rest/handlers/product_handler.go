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

// ProductHandler handles product operations
type ProductHandler struct {
	productService ports.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService ports.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, req bunrouter.Request) error {
	var createReq domain.CreateProductRequest
	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	product, err := h.productService.CreateProduct(req.Context(), &createReq)
	if err != nil {
		http.Error(w, "Failed to create product: "+err.Error(), http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(product)
}

// GetProduct retrieves a product by ID
func (h *ProductHandler) GetProduct(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return err
	}

	product, err := h.productService.GetProduct(req.Context(), id)
	if err != nil {
		http.Error(w, "Product not found: "+err.Error(), http.StatusNotFound)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(product)
}

// GetProducts retrieves all products with pagination and filtering
func (h *ProductHandler) GetProducts(w http.ResponseWriter, req bunrouter.Request) error {
	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")
	categoryIDStr := req.URL.Query().Get("category_id")
	isActiveStr := req.URL.Query().Get("is_active")

	limit := 10
	offset := 0
	var categoryID *uuid.UUID
	var isActive *bool

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

	if categoryIDStr != "" {
		if id, err := uuid.Parse(categoryIDStr); err == nil {
			categoryID = &id
		}
	}

	if isActiveStr != "" {
		if active, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &active
		}
	}

	var products []*domain.Product
	var err error

	if categoryID != nil {
		products, err = h.productService.GetProductsByCategory(req.Context(), *categoryID, limit, offset)
	} else if isActive != nil && *isActive {
		products, err = h.productService.GetActiveProducts(req.Context(), limit, offset)
	} else {
		products, err = h.productService.GetProducts(req.Context(), limit, offset)
	}
	if err != nil {
		http.Error(w, "Failed to get products: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	response := map[string]interface{}{
		"products":    products,
		"limit":       limit,
		"offset":      offset,
		"category_id": categoryID,
		"is_active":   isActive,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// UpdateProduct updates an existing product
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return err
	}

	var updateReq domain.UpdateProductRequest
	if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	product, err := h.productService.UpdateProduct(req.Context(), id, &updateReq)
	if err != nil {
		http.Error(w, "Failed to update product: "+err.Error(), http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(product)
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return err
	}

	err = h.productService.DeleteProduct(req.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete product: "+err.Error(), http.StatusNotFound)
		return err
	}

	response := map[string]string{
		"message": "Product deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registers product routes
func (h *ProductHandler) RegisterRoutes(router *bunrouter.Router) {
	api := router.NewGroup("/api/products")
	api.POST("", h.CreateProduct)
	api.GET("/:id", h.GetProduct)
	api.GET("", h.GetProducts)
	api.PUT("/:id", h.UpdateProduct)
	api.DELETE("/:id", h.DeleteProduct)
}
