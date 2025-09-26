package graphql

import (
	"context"
	"fmt"
	"net/http"

	"silbackendassessment/internal/adapters/middleware"
	graphpkg "silbackendassessment/internal/api/graphql/graph"
	modelspkg "silbackendassessment/internal/api/graphql/graph/model"
	resolverspkg "silbackendassessment/internal/api/graphql/resolvers"
	"silbackendassessment/internal/core/ports"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/uptrace/bunrouter"
)

// RouterConfig holds the configuration for the GraphQL router
type RouterConfig struct {
	UserService         ports.UserService
	CustomerService     ports.CustomerService
	CategoryService     ports.CategoryService
	ProductService      ports.ProductService
	OrderService        ports.OrderService
	NotificationService ports.NotificationService
	AuthMiddleware      *middleware.AuthMiddleware
}

// NewRouter creates a new GraphQL router with all handlers registered
func NewRouter(config *RouterConfig) *bunrouter.Router {
	router := bunrouter.New()

	// Initialize gqlgen server with DI resolver
	r := resolverspkg.NewResolver(
		config.UserService,
		config.CustomerService,
		config.CategoryService,
		config.ProductService,
		config.OrderService,
		config.NotificationService,
	)

	directives := graphpkg.DirectiveRoot{
		Auth: func(ctx context.Context, obj interface{}, next graphql.Resolver, scope *modelspkg.AuthScope) (res interface{}, err error) {
			// Default to ANY if not specified
			s := modelspkg.AuthScopeAny
			if scope != nil {
				s = *scope
			}
			// ANY: allow if either user or customer present
			// USER: require user context
			// CUSTOMER: require customer context
			if s == modelspkg.AuthScopeAny {
				if _, ok := middleware.GetUserFromContext(ctx); ok {
					return next(ctx)
				}
				if _, ok := middleware.GetCustomerFromContext(ctx); ok {
					return next(ctx)
				}
				return nil, fmt.Errorf("unauthorized")
			}
			if s == modelspkg.AuthScopeUser {
				if _, ok := middleware.GetUserFromContext(ctx); ok {
					return next(ctx)
				}
				return nil, fmt.Errorf("user authorization required")
			}
			if s == modelspkg.AuthScopeCustomer {
				if _, ok := middleware.GetCustomerFromContext(ctx); ok {
					return next(ctx)
				}
				return nil, fmt.Errorf("customer authorization required")
			}
			return next(ctx)
		},
	}

	schema := graphpkg.NewExecutableSchema(graphpkg.Config{Resolvers: r, Directives: directives})
	srv := handler.NewDefaultServer(schema)

	// GraphQL endpoint (mounted at /graphql)
	postHandler := func(w http.ResponseWriter, req bunrouter.Request) error {
		srv.ServeHTTP(w, req.Request)
		return nil
	}
	if config.AuthMiddleware != nil {
		postHandler = config.AuthMiddleware.RequireCustomerAuth(postHandler)
	}
	router.POST("/", postHandler)

	// GraphQL Playground (public)
	router.GET("/playground", func(w http.ResponseWriter, req bunrouter.Request) error {
		playground.Handler("GraphQL", "/graphql").ServeHTTP(w, req.Request)
		return nil
	})

	// GraphQL schema endpoint (simple ok)
	router.GET("/schema", func(w http.ResponseWriter, req bunrouter.Request) error {
		w.Header().Set("Content-Type", "application/json")
		return bunrouter.JSON(w, map[string]string{"status": "ok"})
	})

	return router
}
