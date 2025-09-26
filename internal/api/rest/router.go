package rest

import (
	"net/http"

	"silbackendassessment/internal/adapters/middleware"
	"silbackendassessment/internal/api/rest/handlers"
	"silbackendassessment/internal/core/ports"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

// RouterConfig holds the configuration for the REST router
type RouterConfig struct {
	AuthMiddleware      *middleware.AuthMiddleware
	UserService         ports.UserService
	CustomerService     ports.CustomerService
	CategoryService     ports.CategoryService
	ProductService      ports.ProductService
	OrderService        ports.OrderService
	NotificationService ports.NotificationService
	AuthService         ports.AuthService
}

// NewRouter creates a new REST router with all handlers registered
func NewRouter(config *RouterConfig) *bunrouter.Router {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(config.UserService)
	customerAuthHandler := handlers.NewCustomerAuthHandler(config.CustomerService, config.AuthService)
	categoryHandler := handlers.NewCategoryHandler(config.CategoryService)
	productHandler := handlers.NewProductHandler(config.ProductService)
	orderHandler := handlers.NewOrderHandler(config.OrderService)
	notificationHandler := handlers.NewNotificationHandler(config.NotificationService)

	// Health check endpoint
	router.GET("/health", func(w http.ResponseWriter, req bunrouter.Request) error {
		return bunrouter.JSON(w, map[string]string{
			"status":  "ok",
			"service": "sil-backend-assessment",
		})
	})

	// Register all routes
	userHandler.RegisterRoutes(router)
	customerAuthHandler.RegisterRoutes(router, config.AuthMiddleware)
	categoryHandler.RegisterRoutes(router)
	productHandler.RegisterRoutes(router)
	orderHandler.RegisterRoutes(router)
	notificationHandler.RegisterRoutes(router)

	return router
}
