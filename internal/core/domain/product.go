package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Product represents a product in the system
type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID          uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name        string    `bun:"name,notnull" json:"name"`
	Description string    `bun:"description" json:"description"`
	SKU         string    `bun:"sku,unique,notnull" json:"sku"`
	Price       float64   `bun:"price,notnull" json:"price"`
	Stock       int       `bun:"stock,notnull,default:0" json:"stock"`
	CategoryID  uuid.UUID `bun:"category_id,type:uuid,notnull" json:"category_id"`
	IsActive    bool      `bun:"is_active,notnull,default:true" json:"is_active"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	// Relations
	Category   Category    `bun:"rel:belongs-to,join:category_id=id" json:"category"`
	OrderItems []OrderItem `bun:"rel:has-many,join:id=product_id" json:"order_items,omitempty"`
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	SKU         string    `json:"sku" validate:"required"`
	Price       float64   `json:"price" validate:"required,min=0"`
	Stock       int       `json:"stock" validate:"min=0"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	IsActive    bool      `json:"is_active"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	SKU         *string    `json:"sku,omitempty"`
	Price       *float64   `json:"price,omitempty"`
	Stock       *int       `json:"stock,omitempty"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
}
