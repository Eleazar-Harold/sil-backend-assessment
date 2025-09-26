package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/uptrace/bunrouter"
	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"
)

// CategoryHandler handles category operations
type CategoryHandler struct {
	categoryService ports.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService ports.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, req bunrouter.Request) error {
	var createReq domain.CreateCategoryRequest
	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	category, err := h.categoryService.CreateCategory(req.Context(), &createReq)
	if err != nil {
		http.Error(w, "Failed to create category: "+err.Error(), http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(category)
}

// GetCategory retrieves a category by ID
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return err
	}

	category, err := h.categoryService.GetCategory(req.Context(), id)
	if err != nil {
		http.Error(w, "Category not found: "+err.Error(), http.StatusNotFound)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(category)
}

// GetCategories retrieves all categories with pagination
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, req bunrouter.Request) error {
	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")

	limit := 10
	offset := 0

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

	categories, err := h.categoryService.GetCategories(req.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to get categories: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	response := map[string]interface{}{
		"categories": categories,
		"limit":      limit,
		"offset":     offset,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// UpdateCategory updates an existing category
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return err
	}

	var updateReq domain.UpdateCategoryRequest
	if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	category, err := h.categoryService.UpdateCategory(req.Context(), id, &updateReq)
	if err != nil {
		http.Error(w, "Failed to update category: "+err.Error(), http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(category)
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return err
	}

	err = h.categoryService.DeleteCategory(req.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete category: "+err.Error(), http.StatusNotFound)
		return err
	}

	response := map[string]string{
		"message": "Category deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registers category routes
func (h *CategoryHandler) RegisterRoutes(router *bunrouter.Router) {
	api := router.NewGroup("/api/categories")
	api.POST("", h.CreateCategory)
	api.GET("/:id", h.GetCategory)
	api.GET("", h.GetCategories)
	api.PUT("/:id", h.UpdateCategory)
	api.DELETE("/:id", h.DeleteCategory)
}

