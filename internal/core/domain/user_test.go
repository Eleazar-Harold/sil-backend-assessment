package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Run("Create user with valid data", func(t *testing.T) {
		user := &User{
			ID:        uuid.New(),
			Name:      "John Doe",
			Email:     "john.doe@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, user)
		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john.doe@example.com", user.Email)
		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())
	})

	t.Run("Create user with empty fields", func(t *testing.T) {
		user := &User{
			ID:        uuid.New(),
			Name:      "",
			Email:     "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, user)
		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.Equal(t, "", user.Name)
		assert.Equal(t, "", user.Email)
	})

	t.Run("Create user with special characters", func(t *testing.T) {
		user := &User{
			ID:        uuid.New(),
			Name:      "José María O'Connor-Smith",
			Email:     "josé.maría+test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		assert.NotNil(t, user)
		assert.Equal(t, "José María O'Connor-Smith", user.Name)
		assert.Equal(t, "josé.maría+test@example.com", user.Email)
	})
}

func TestCreateUserRequest(t *testing.T) {
	t.Run("Create user request with valid data", func(t *testing.T) {
		req := &CreateUserRequest{
			Name:  "Jane Doe",
			Email: "jane.doe@example.com",
		}

		assert.NotNil(t, req)
		assert.Equal(t, "Jane Doe", req.Name)
		assert.Equal(t, "jane.doe@example.com", req.Email)
	})

	t.Run("Create user request with empty fields", func(t *testing.T) {
		req := &CreateUserRequest{
			Name:  "",
			Email: "",
		}

		assert.NotNil(t, req)
		assert.Equal(t, "", req.Name)
		assert.Equal(t, "", req.Email)
	})

	t.Run("Create user request with special characters", func(t *testing.T) {
		req := &CreateUserRequest{
			Name:  "François Müller",
			Email: "françois.müller+tag@example.com",
		}

		assert.NotNil(t, req)
		assert.Equal(t, "François Müller", req.Name)
		assert.Equal(t, "françois.müller+tag@example.com", req.Email)
	})
}

func TestUpdateUserRequest(t *testing.T) {
	t.Run("Update user request with all fields", func(t *testing.T) {
		newName := "Updated Name"
		newEmail := "updated@example.com"

		req := &UpdateUserRequest{
			Name:  &newName,
			Email: &newEmail,
		}

		assert.NotNil(t, req)
		assert.NotNil(t, req.Name)
		assert.NotNil(t, req.Email)
		assert.Equal(t, "Updated Name", *req.Name)
		assert.Equal(t, "updated@example.com", *req.Email)
	})

	t.Run("Update user request with only name", func(t *testing.T) {
		newName := "Updated Name Only"

		req := &UpdateUserRequest{
			Name:  &newName,
			Email: nil,
		}

		assert.NotNil(t, req)
		assert.NotNil(t, req.Name)
		assert.Nil(t, req.Email)
		assert.Equal(t, "Updated Name Only", *req.Name)
	})

	t.Run("Update user request with only email", func(t *testing.T) {
		newEmail := "updated-email-only@example.com"

		req := &UpdateUserRequest{
			Name:  nil,
			Email: &newEmail,
		}

		assert.NotNil(t, req)
		assert.Nil(t, req.Name)
		assert.NotNil(t, req.Email)
		assert.Equal(t, "updated-email-only@example.com", *req.Email)
	})

	t.Run("Update user request with no fields", func(t *testing.T) {
		req := &UpdateUserRequest{
			Name:  nil,
			Email: nil,
		}

		assert.NotNil(t, req)
		assert.Nil(t, req.Name)
		assert.Nil(t, req.Email)
	})

	t.Run("Update user request with empty string values", func(t *testing.T) {
		emptyName := ""
		emptyEmail := ""

		req := &UpdateUserRequest{
			Name:  &emptyName,
			Email: &emptyEmail,
		}

		assert.NotNil(t, req)
		assert.NotNil(t, req.Name)
		assert.NotNil(t, req.Email)
		assert.Equal(t, "", *req.Name)
		assert.Equal(t, "", *req.Email)
	})
}
