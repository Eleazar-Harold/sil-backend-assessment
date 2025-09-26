package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"

	"silbackendassessment/internal/adapters/auth"
	"silbackendassessment/internal/adapters/middleware"
	"silbackendassessment/internal/adapters/notifications"
	"silbackendassessment/internal/adapters/repositories"
	"silbackendassessment/internal/api/graphql"
	"silbackendassessment/internal/api/rest"
	"silbackendassessment/internal/api/rest/handlers"
	"silbackendassessment/internal/config"
	"silbackendassessment/internal/core/oidc"
	"silbackendassessment/internal/core/services"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	orderItemRepo := repositories.NewOrderItemRepository(db)

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.RefreshSecret,
		cfg.Auth.JWTExpiry,
		cfg.Auth.RefreshExpiry,
	)

	// Initialize OIDC provider if enabled
	var oidcProvider oidc.Provider
	if cfg.OIDC.Enabled {
		oidcProvider, err = auth.NewOIDCProvider(
			cfg.OIDC.ProviderURL,
			cfg.OIDC.ClientID,
			cfg.OIDC.ClientSecret,
			cfg.OIDC.RedirectURL,
			cfg.OIDC.Scopes,
		)
		if err != nil {
			log.Fatalf("Failed to initialize OIDC provider: %v", err)
		}
	}

	// Initialize notification adapters
	emailClient := notifications.NewEmailClient(&notifications.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
		TLS:      cfg.SMTP.TLS,
	})

	smsClient := notifications.NewSMSClient(&notifications.ATConfig{
		APIKey:   cfg.AT.APIKey,
		Username: cfg.AT.Username,
		BaseURL:  cfg.AT.BaseURL,
	})

	// Initialize notification service
	notificationService := services.NewNotificationService(emailClient, smsClient)

	// Initialize services
	authService := services.NewAuthService(userRepo, customerRepo, jwtManager, oidcProvider)
	userService := services.NewUserService(userRepo)
	customerService := services.NewCustomerService(customerRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	orderService := services.NewOrderService(orderRepo, orderItemRepo, customerRepo, productRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	customerHandler := handlers.NewCustomerAuthHandler(customerService, authService)
	oidcHandler := handlers.NewOIDCHandler(authService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup main router
	router := bunrouter.New(
		bunrouter.WithMiddleware(reqlog.NewMiddleware()),
		bunrouter.WithMiddleware(corsMiddleware()),
	)

	// Initialize REST API router
	restConfig := &rest.RouterConfig{
		AuthMiddleware:      authMiddleware,
		UserService:         userService,
		CustomerService:     customerService,
		CategoryService:     categoryService,
		ProductService:      productService,
		OrderService:        orderService,
		NotificationService: notificationService,
		AuthService:         authService,
	}
	restRouter := rest.NewRouter(restConfig)

	// Initialize GraphQL router
	graphqlConfig := &graphql.RouterConfig{
		UserService:         userService,
		CustomerService:     customerService,
		CategoryService:     categoryService,
		ProductService:      productService,
		OrderService:        orderService,
		NotificationService: notificationService,
		AuthMiddleware:      authMiddleware,
	}
	graphqlRouter := graphql.NewRouter(graphqlConfig)

	// Register routes
	registerRoutes(router, userHandler, customerHandler, oidcHandler, authMiddleware, restRouter, graphqlRouter)

	// Start server
	log.Printf("Server starting on port %d", cfg.Server.RESTPort)
	log.Printf("REST API available at: http://localhost:%d/api", cfg.Server.RESTPort)
	log.Printf("GraphQL API available at: http://localhost:%d/graphql", cfg.Server.RESTPort)
	log.Printf("GraphQL Playground available at: http://localhost:%d/graphql/playground", cfg.Server.RESTPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.RESTPort), router))
}

func connectDB(cfg *config.Config) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.GetDSN())))

	// Set connection limits
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(time.Hour)

	// Check connection
	if err := sqldb.Ping(); err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	return db, nil
}

func registerRoutes(
	router *bunrouter.Router,
	userHandler *handlers.UserHandler,
	customerHandler *handlers.CustomerAuthHandler,
	oidcHandler *handlers.OIDCHandler,
	authMiddleware *middleware.AuthMiddleware,
	restRouter *bunrouter.Router,
	graphqlRouter *bunrouter.Router,
) {
	// Public routes
	router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		log.Printf("Root hit: %s %s", req.Method, req.URL.Path)
		return bunrouter.JSON(w, map[string]interface{}{
			"message": "SIL Backend Assessment API",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"rest":    "/api",
				"graphql": "/graphql",
				"docs":    "/docs",
				"health":  "/api/health",
			},
		})
	})

	// Health check endpoint (accessible at /api/health)
	router.GET("/api/health", func(w http.ResponseWriter, req bunrouter.Request) error {
		return bunrouter.JSON(w, map[string]string{
			"status":  "ok",
			"service": "sil-backend-assessment",
			"version": "1.0.0",
		})
	})

	// Mount REST API router
	router.WithGroup("/api", func(g *bunrouter.Group) {
		forward := func(w http.ResponseWriter, req bunrouter.Request) error {
			// Forward without stripping to preserve /api prefix expected by restRouter
			r := req.Request.Clone(req.Context())
			log.Printf("REST forward: %s %s -> %s", req.Method, req.URL.Path, r.URL.Path)
			restRouter.ServeHTTP(w, r)
			return nil
		}
		g.GET("/*path", forward)
		g.POST("/*path", forward)
		g.PUT("/*path", forward)
		g.DELETE("/*path", forward)
		g.OPTIONS("/*path", forward)
	})

	// Mount GraphQL router
	router.WithGroup("/graphql", func(g *bunrouter.Group) {
		forward := func(w http.ResponseWriter, req bunrouter.Request) error {
			// Strip group prefix and forward to GraphQL router
			r := req.Request.Clone(req.Context())
			if strings.HasPrefix(r.URL.Path, "/graphql") {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, "/graphql")
				if r.URL.Path == "" {
					r.URL.Path = "/"
				}
			}
			graphqlRouter.ServeHTTP(w, r)
			return nil
		}
		// Ensure root of group matches exactly for /graphql
		g.POST("/", forward)
		g.GET("/*path", forward)
		g.PUT("/*path", forward)
		g.DELETE("/*path", forward)
		g.OPTIONS("/*path", forward)
	})

	// OIDC authentication routes
	oidcHandler.RegisterRoutes(router)

	// User routes (traditional JWT auth) - handled by REST router

	// Customer routes (OIDC auth)
	customerHandler.RegisterRoutes(router, authMiddleware)

	// API documentation endpoints
	router.GET("/docs", func(w http.ResponseWriter, req bunrouter.Request) error {
		return bunrouter.JSON(w, map[string]interface{}{
			"api_documentation": map[string]string{
				"rest":               "http://localhost:8080/api",
				"graphql":            "http://localhost:8080/graphql",
				"graphql_playground": "http://localhost:8080/graphql/playground",
				"health":             "http://localhost:8080/api/health",
			},
			"documentation_files": map[string]string{
				"main_api_docs":        "API_DOCUMENTATION.md",
				"quick_reference":      "ENDPOINTS_QUICK_REFERENCE.md",
				"graphql_docs":         "GRAPHQL_API_DOCUMENTATION.md",
				"graphql_complete_ref": "GRAPHQL_API_COMPLETE_REFERENCE.md",
				"notification_docs":    "NOTIFICATION_API_DOCUMENTATION.md",
				"oidc_setup":           "OIDC_SETUP.md",
			},
			"authentication": map[string]interface{}{
				"jwt_token_format":  "Authorization: Bearer <jwt-token>",
				"oidc_token_format": "Authorization: Bearer <oidc-token>",
				"auth_scopes": map[string]string{
					"ANY":      "Either JWT or OIDC token",
					"USER":     "JWT token required",
					"CUSTOMER": "OIDC token required",
				},
			},
		})
	})

	// Swagger/OpenAPI documentation placeholder
	router.GET("/swagger.json", func(w http.ResponseWriter, req bunrouter.Request) error {
		w.Header().Set("Content-Type", "application/json")
		return bunrouter.JSON(w, map[string]interface{}{
			"openapi": "3.0.0",
			"info": map[string]string{
				"title":       "SIL Backend Assessment API",
				"version":     "1.0.0",
				"description": "REST API for SIL Backend Assessment with GraphQL support",
			},
			"servers": []map[string]string{
				{"url": "http://localhost:8080/api"},
			},
			"message": "Full API documentation available in markdown files in the repository",
		})
	})

	// Redoc documentation placeholder
	router.GET("/redoc.html", func(w http.ResponseWriter, req bunrouter.Request) error {
		w.Header().Set("Content-Type", "text/html")
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>SIL Backend Assessment API Documentation</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
        .endpoint { background: #f4f4f4; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .method { font-weight: bold; color: #2c3e50; }
        .path { font-family: monospace; background: #e8e8e8; padding: 2px 5px; }
    </style>
</head>
<body>
    <h1>SIL Backend Assessment API Documentation</h1>
    <p>This is a placeholder for Redoc documentation. Full documentation is available in the following markdown files:</p>
    
    <ul>
        <li><strong>API_DOCUMENTATION.md</strong> - Complete REST and GraphQL API documentation</li>
        <li><strong>ENDPOINTS_QUICK_REFERENCE.md</strong> - Quick reference for all endpoints</li>
        <li><strong>GRAPHQL_API_DOCUMENTATION.md</strong> - Detailed GraphQL documentation</li>
        <li><strong>GRAPHQL_API_COMPLETE_REFERENCE.md</strong> - Complete GraphQL reference</li>
        <li><strong>NOTIFICATION_API_DOCUMENTATION.md</strong> - Email and SMS notification API</li>
        <li><strong>OIDC_SETUP.md</strong> - OIDC authentication setup guide</li>
    </ul>

    <h2>Quick Start</h2>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/api/health</span> - Health check
    </div>
    <div class="endpoint">
        <span class="method">POST</span> <span class="path">/graphql</span> - GraphQL endpoint
    </div>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/graphql/playground</span> - GraphQL playground
    </div>
    <div class="endpoint">
        <span class="method">GET</span> <span class="path">/docs</span> - API documentation index
    </div>

    <h2>Authentication</h2>
    <p>Use JWT or OIDC tokens in the Authorization header:</p>
    <pre>Authorization: Bearer &lt;token&gt;</pre>
</body>
</html>`
		w.Write([]byte(html))
		return nil
	})
}

func corsMiddleware() bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if req.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return nil
			}

			return next(w, req)
		}
	}
}
