package graph

import (
	"silbackendassessment/internal/core/ports"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userService     ports.UserService
	customerService ports.CustomerService
	categoryService ports.CategoryService
	productService  ports.ProductService
	orderService    ports.OrderService
}

// NewResolver creates a new resolver with all required services
func NewResolver(
	userService ports.UserService,
	customerService ports.CustomerService,
	categoryService ports.CategoryService,
	productService ports.ProductService,
	orderService ports.OrderService,
) *Resolver {
	return &Resolver{
		userService:     userService,
		customerService: customerService,
		categoryService: categoryService,
		productService:  productService,
		orderService:    orderService,
	}
}
