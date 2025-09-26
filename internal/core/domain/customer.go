package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Customer represents a customer in the system
type Customer struct {
	bun.BaseModel `bun:"table:customers,alias:c"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	FirstName string    `bun:"first_name,notnull" json:"first_name"`
	LastName  string    `bun:"last_name,notnull" json:"last_name"`
	Email     string    `bun:"email,unique,notnull" json:"email"`
	Phone     string    `bun:"phone" json:"phone"`
	Address   string    `bun:"address" json:"address"`
	City      string    `bun:"city" json:"city"`
	State     string    `bun:"state" json:"state"`
	ZipCode   string    `bun:"zip_code" json:"zip_code"`
	Country   string    `bun:"country" json:"country"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	// Relations
	Orders []Order `bun:"rel:has-many,join:id=customer_id" json:"orders,omitempty"`
}

// CreateCustomerRequest represents the request to create a customer
type CreateCustomerRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
}

// UpdateCustomerRequest represents the request to update a customer
type UpdateCustomerRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Address   *string `json:"address,omitempty"`
	City      *string `json:"city,omitempty"`
	State     *string `json:"state,omitempty"`
	ZipCode   *string `json:"zip_code,omitempty"`
	Country   *string `json:"country,omitempty"`
}
