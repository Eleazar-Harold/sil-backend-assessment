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

type UserHandler struct {
	userService ports.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, req bunrouter.Request) error {
	var createReq domain.CreateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	user, err := h.userService.CreateUser(req.Context(), &createReq)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}

	user, err := h.userService.GetUser(req.Context(), id)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, req bunrouter.Request) error {
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

	users, err := h.userService.GetUsers(req.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to get users: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	response := map[string]interface{}{
		"users":  users,
		"limit":  limit,
		"offset": offset,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}

	var updateReq domain.UpdateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	user, err := h.userService.UpdateUser(req.Context(), id, &updateReq)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusBadRequest)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, req bunrouter.Request) error {
	idStr := req.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}

	if err := h.userService.DeleteUser(req.Context(), id); err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusNotFound)
		return err
	}

	response := map[string]string{
		"message": "User deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registers user routes
func (h *UserHandler) RegisterRoutes(router *bunrouter.Router) {
	api := router.NewGroup("/api/users")
	api.POST("", h.CreateUser)
	api.GET("/:id", h.GetUser)
	api.GET("", h.GetUsers)
	api.PUT("/:id", h.UpdateUser)
	api.DELETE("/:id", h.DeleteUser)
}
