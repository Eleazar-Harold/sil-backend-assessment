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

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := testutils.NewMockUserRepository()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Create user successfully", func(t *testing.T) {
		req := &domain.CreateUserRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		user, err := service.CreateUser(ctx, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if user == nil {
			t.Error("Expected user to be created")
		}

		if user.Name != "John Doe" {
			t.Errorf("Expected name to be 'John Doe', got: %s", user.Name)
		}

		if user.Email != "john@example.com" {
			t.Errorf("Expected email to be 'john@example.com', got: %s", user.Email)
		}

		if user.ID == uuid.Nil {
			t.Error("Expected user ID to be set")
		}

		if user.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		if user.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("Create user with existing email", func(t *testing.T) {
		// Set up existing user
		existingUser := &domain.User{
			ID:        uuid.New(),
			Name:      "Existing User",
			Email:     "existing@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.UsersByEmail["existing@example.com"] = existingUser

		req := &domain.CreateUserRequest{
			Name:  "New User",
			Email: "existing@example.com",
		}

		user, err := service.CreateUser(ctx, req)

		if err == nil {
			t.Error("Expected error for existing email")
		}

		if user != nil {
			t.Error("Expected user to be nil")
		}

		if err == nil || err.Error() != "user with email existing@example.com already exists" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("Repository error during creation", func(t *testing.T) {
		mockRepo.CreateError = errors.New("database error")

		req := &domain.CreateUserRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		user, err := service.CreateUser(ctx, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if user != nil {
			t.Error("Expected user to be nil")
		}
	})
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := testutils.NewMockUserRepository()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Get existing user", func(t *testing.T) {
		userID := uuid.New()
		expectedUser := &domain.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Users[userID] = expectedUser

		user, err := service.GetUser(ctx, userID)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if user == nil {
			t.Error("Expected user to be returned")
		}

		if user.ID != userID {
			t.Errorf("Expected user ID %s, got: %s", userID, user.ID)
		}
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		userID := uuid.New()

		user, err := service.GetUser(ctx, userID)

		if err == nil {
			t.Error("Expected error for non-existent user")
		}

		if user != nil {
			t.Error("Expected user to be nil")
		}

		if !errors.Is(err, testutils.ErrUserNotFound) {
			t.Errorf("Expected ErrUserNotFound, got: %v", err)
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		userID := uuid.New()
		mockRepo.GetByIDError = errors.New("database error")

		user, err := service.GetUser(ctx, userID)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if user != nil {
			t.Error("Expected user to be nil")
		}
	})
}

func TestUserService_GetUsers(t *testing.T) {
	mockRepo := testutils.NewMockUserRepository()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Get users successfully", func(t *testing.T) {
		users := []*domain.User{
			{
				ID:        uuid.New(),
				Name:      "User 1",
				Email:     "user1@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				Name:      "User 2",
				Email:     "user2@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRepo.AllUsers = users

		result, err := service.GetUsers(ctx, 10, 0)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 users, got: %d", len(result))
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo.GetAllError = errors.New("database error")

		users, err := service.GetUsers(ctx, 10, 0)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if users != nil {
			t.Error("Expected users to be nil")
		}
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := testutils.NewMockUserRepository()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Update user successfully", func(t *testing.T) {
		userID := uuid.New()
		existingUser := &domain.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Users[userID] = existingUser

		newName := "John Updated"
		req := &domain.UpdateUserRequest{
			Name: &newName,
		}

		user, err := service.UpdateUser(ctx, userID, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if user == nil {
			t.Error("Expected user to be returned")
		}

		if user.Name != "John Updated" {
			t.Errorf("Expected name to be 'John Updated', got: %s", user.Name)
		}

		if user.Email != "john@example.com" {
			t.Errorf("Expected email to remain unchanged, got: %s", user.Email)
		}
	})

	t.Run("Update non-existent user", func(t *testing.T) {
		userID := uuid.New()
		newName := "John Updated"
		req := &domain.UpdateUserRequest{
			Name: &newName,
		}

		user, err := service.UpdateUser(ctx, userID, req)

		if err == nil {
			t.Error("Expected error for non-existent user")
		}

		if user != nil {
			t.Error("Expected user to be nil")
		}

		if !errors.Is(err, testutils.ErrUserNotFound) {
			t.Errorf("Expected ErrUserNotFound, got: %v", err)
		}
	})

	t.Run("Update user email", func(t *testing.T) {
		userID := uuid.New()
		existingUser := &domain.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Users[userID] = existingUser

		newEmail := "john.updated@example.com"
		req := &domain.UpdateUserRequest{
			Email: &newEmail,
		}

		user, err := service.UpdateUser(ctx, userID, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if user.Email != "john.updated@example.com" {
			t.Errorf("Expected email to be 'john.updated@example.com', got: %s", user.Email)
		}
	})

	t.Run("Repository error during update", func(t *testing.T) {
		userID := uuid.New()
		existingUser := &domain.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Users[userID] = existingUser
		mockRepo.UpdateError = errors.New("database error")

		newName := "John Updated"
		req := &domain.UpdateUserRequest{
			Name: &newName,
		}

		user, err := service.UpdateUser(ctx, userID, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if user != nil {
			t.Error("Expected user to be nil")
		}
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := testutils.NewMockUserRepository()
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Delete user successfully", func(t *testing.T) {
		userID := uuid.New()
		existingUser := &domain.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Users[userID] = existingUser

		err := service.DeleteUser(ctx, userID)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify user was deleted
		if _, exists := mockRepo.Users[userID]; exists {
			t.Error("Expected user to be deleted")
		}
	})

	t.Run("Delete non-existent user", func(t *testing.T) {
		userID := uuid.New()

		err := service.DeleteUser(ctx, userID)

		if err == nil {
			t.Error("Expected error for non-existent user")
		}

		if !errors.Is(err, testutils.ErrUserNotFound) {
			t.Errorf("Expected ErrUserNotFound, got: %v", err)
		}
	})

	t.Run("Repository error during deletion", func(t *testing.T) {
		userID := uuid.New()
		existingUser := &domain.User{
			ID:        userID,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Users[userID] = existingUser
		mockRepo.DeleteError = errors.New("database error")

		err := service.DeleteUser(ctx, userID)

		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}
