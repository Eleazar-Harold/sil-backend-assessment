package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"silbackendassessment/internal/core/ports"
)

type Resolver struct {
	userService         ports.UserService
	customerService     ports.CustomerService
	categoryService     ports.CategoryService
	productService      ports.ProductService
	orderService        ports.OrderService
	notificationService ports.NotificationService
}

func NewResolver(
	userService ports.UserService,
	customerService ports.CustomerService,
	categoryService ports.CategoryService,
	productService ports.ProductService,
	orderService ports.OrderService,
	notificationService ports.NotificationService,
) *Resolver {
	return &Resolver{
		userService:         userService,
		customerService:     customerService,
		categoryService:     categoryService,
		productService:      productService,
		orderService:        orderService,
		notificationService: notificationService,
	}
}
