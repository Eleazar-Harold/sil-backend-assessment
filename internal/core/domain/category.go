package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Category represents a product category
type Category struct {
	bun.BaseModel `bun:"table:categories,alias:cat"`

	ID          uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name        string     `bun:"name,unique,notnull" json:"name"`
	Description string     `bun:"description" json:"description"`
	ParentID    *uuid.UUID `bun:"parent_id,type:uuid" json:"parent_id"`
	CreatedAt   time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	// Relations
	Parent   *Category  `bun:"rel:belongs-to,join:parent_id=id" json:"parent,omitempty"`
	Children []Category `bun:"rel:has-many,join:id=parent_id" json:"children,omitempty"`
	Products []Product  `bun:"rel:has-many,join:id=category_id" json:"products,omitempty"`
}

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
}
