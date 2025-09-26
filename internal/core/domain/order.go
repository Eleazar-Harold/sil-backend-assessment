package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// Order represents an order in the system
type Order struct {
	bun.BaseModel `bun:"table:orders,alias:o"`

	ID              uuid.UUID   `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	CustomerID      uuid.UUID   `bun:"customer_id,type:uuid,notnull" json:"customer_id"`
	OrderNumber     string      `bun:"order_number,unique,notnull" json:"order_number"`
	Status          OrderStatus `bun:"status,notnull,default:'pending'" json:"status"`
	TotalAmount     float64     `bun:"total_amount,notnull" json:"total_amount"`
	ShippingAddress string      `bun:"shipping_address,notnull" json:"shipping_address"`
	BillingAddress  string      `bun:"billing_address,notnull" json:"billing_address"`
	Notes           string      `bun:"notes" json:"notes"`
	OrderDate       time.Time   `bun:"order_date,nullzero,notnull,default:current_timestamp" json:"order_date"`
	ShippedDate     *time.Time  `bun:"shipped_date" json:"shipped_date"`
	DeliveredDate   *time.Time  `bun:"delivered_date" json:"delivered_date"`
	CreatedAt       time.Time   `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time   `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	// Relations
	Customer   Customer    `bun:"rel:belongs-to,join:customer_id=id" json:"customer"`
	OrderItems []OrderItem `bun:"rel:has-many,join:id=order_id" json:"order_items"`
}

// OrderItem represents an item within an order
type OrderItem struct {
	bun.BaseModel `bun:"table:order_items,alias:oi"`

	ID         uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	OrderID    uuid.UUID `bun:"order_id,type:uuid,notnull" json:"order_id"`
	ProductID  uuid.UUID `bun:"product_id,type:uuid,notnull" json:"product_id"`
	Quantity   int       `bun:"quantity,notnull" json:"quantity"`
	UnitPrice  float64   `bun:"unit_price,notnull" json:"unit_price"`
	TotalPrice float64   `bun:"total_price,notnull" json:"total_price"`
	CreatedAt  time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	// Relations
	Order   Order   `bun:"rel:belongs-to,join:order_id=id" json:"order"`
	Product Product `bun:"rel:belongs-to,join:product_id=id" json:"product"`
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	CustomerID      uuid.UUID                `json:"customer_id" validate:"required"`
	ShippingAddress string                   `json:"shipping_address" validate:"required"`
	BillingAddress  string                   `json:"billing_address" validate:"required"`
	Notes           string                   `json:"notes"`
	OrderItems      []CreateOrderItemRequest `json:"order_items" validate:"required,min=1"`
}

// CreateOrderItemRequest represents the request to create an order item
type CreateOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

// UpdateOrderRequest represents the request to update an order
type UpdateOrderRequest struct {
	Status          *OrderStatus `json:"status,omitempty"`
	ShippingAddress *string      `json:"shipping_address,omitempty"`
	BillingAddress  *string      `json:"billing_address,omitempty"`
	Notes           *string      `json:"notes,omitempty"`
	ShippedDate     *time.Time   `json:"shipped_date,omitempty"`
	DeliveredDate   *time.Time   `json:"delivered_date,omitempty"`
}
