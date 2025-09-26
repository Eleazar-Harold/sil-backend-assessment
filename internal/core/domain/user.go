package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// User represents the core domain model
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	Email     string    `bun:"email,unique,notnull" json:"email"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}