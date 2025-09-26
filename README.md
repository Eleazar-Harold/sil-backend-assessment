# sil_backend_assessment

A Go assessment with implemented graphql and rest endpoints

## Architecture

This project follows the **Adapter Pattern** (Hexagonal Architecture) principles:

- **Core Domain**: Contains business entities and logic (`internal/core/domain`, `internal/core/services`)
- **Ports**: Define interfaces for adapters (`internal/core/ports`)
- **Adapters**: Implement the ports for external concerns
  - **Handlers**: HTTP request handlers (`internal/adapters/handlers`)
  - **Repositories**: Database operations (`internal/adapters/repositories`)

## Technology Stack

- **Go 1.24**
- **[BunRouter](https://github.com/uptrace/bunrouter)**: HTTP router and middleware
- **[Bun](https://github.com/uptrace/bun)**: SQL-first database toolkit
- **PostgreSQL**: Database
- **UUID**: For primary keys
- **GraphQL**: Via [gqlgen](https://gqlgen.com/) for type-safe GraphQL API
- **JWT & OIDC**: Authentication and authorization
- **SMTP**: Email notifications
- **Africa's Talking**: SMS notifications

## Project Structure

```
sil-backend-assessment/
├── cmd/
│   ├── migrate/
│   │   └── main.go                  # Application migrations point
│   │   ├── migrations/              # Database migrations
│   ├── server/
│   │   └── main.go                  # Application entry point
├── go.mod                           # Go module definition
├── .env.example                     # Environment variables example
├── internal/
│   ├── config/                      # Configuration
│   │   └── config.go
│   ├── core/                        # Core business logic
│   │   ├── domain/                  # Domain models
│   │   │   └── user.go
│   │   ├── ports/                   # Interfaces
│   │   │   ├── user_repository.go
│   │   │   └── user_service.go
│   │   └── services/                # Business logic implementation
│   │       └── user_service.go
│   └── adapters/                    # External adapters
│       ├── handlers/                # HTTP handlers
│       │   ├── user_handler.go
│       │   ├── notification_handler.go
│       │   └── ...
│       ├── notifications/           # Email and SMS adapters
│       │   ├── email.go
│       │   └── sms.go
│       └── repositories/            # Database adapters
│           └── user_repository.go
```

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 12+
- Redis 6+ (optional, for caching)
- Docker and Docker Compose (for containerized development)
- Minikube (for Kubernetes deployment)
- kubectl (for Kubernetes management)
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd sil-backend-assessment
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Update the `.env` file with your database configuration.

4. Install dependencies:
```bash
make deps
# or
go mod tidy
```

5. Run database migrations:
```bash
# Create the database first
createdb sil_backend_assessment_db

# Run migrations using the built-in migration tool
make migrate_up
```

6. Run the application:
```bash
make build
./bin/server
# or
go run cmd/server/main.go -config config.yaml
```

The server will start on port 8080 by default.

## API Endpoints

The SIL Backend Assessment provides comprehensive REST and GraphQL APIs for managing users, customers, categories, products, and orders.

### Quick Reference

- **REST API Base**: `http://localhost:8080/api`
- **GraphQL Endpoint**: `http://localhost:8080/graphql`
- **GraphQL Playground**: `http://localhost:8080/graphql/playground`
- **API Documentation**: `http://localhost:8080/docs`

### Authentication

The API supports two authentication methods:
- **JWT Authentication**: For traditional user authentication
- **OIDC Authentication**: For customer authentication via external providers (Google, Microsoft, Auth0, etc.)

### Available Endpoints

#### REST API
- **Users**: `/api/users` - Full CRUD operations for user management
- **Categories**: `/api/categories` - Product category management with hierarchical support
- **Products**: `/api/products` - Product management with inventory tracking
- **Orders**: `/api/orders` - Order management with status tracking
- **Customer Profile**: `/api/customer/profile` - Customer profile management (OIDC authenticated)
- **OIDC Auth**: `/auth/oidc/*` - OpenID Connect authentication endpoints
- **Health Check**: `/api/health` - API health status

#### GraphQL API
- **Queries**: Retrieve users, customers, categories, products, and orders with flexible filtering
- **Mutations**: Create, update, and delete operations for all entities
- **Relationships**: Nested queries with automatic relationship resolution

### Interactive Documentation

The API includes comprehensive interactive documentation:

- **Swagger UI**: `http://localhost:8080/swagger-ui.html` - Interactive REST API testing
- **ReDoc**: `http://localhost:8080/redoc.html` - Beautiful REST API documentation
- **GraphQL Docs**: `http://localhost:8080/graphql-docs.html` - Complete GraphQL documentation
- **GraphQL Playground**: `http://localhost:8080/graphql/playground` - Interactive GraphQL IDE
- **OpenAPI Spec**: `http://localhost:8080/swagger.json` - Machine-readable API specification

### Comprehensive Documentation

For detailed API documentation including:
- Complete endpoint specifications
- Request/response examples
- GraphQL schema and queries
- Authentication flows
- Error handling

**See: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)**

### OIDC Setup

For OpenID Connect authentication setup and configuration:

**See: [OIDC_SETUP.md](./OIDC_SETUP.md)**

### Documentation Files

- **Interactive Docs**: [docs/](./docs/) - Swagger UI, ReDoc, and GraphQL documentation
- **Quick Reference**: [ENDPOINTS_QUICK_REFERENCE.md](./ENDPOINTS_QUICK_REFERENCE.md)
- **Full API Docs**: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

## Example Usage

### REST API Examples

#### Create a user
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

#### Get all users
```bash
curl -H "Authorization: Bearer <jwt-token>" \
  "http://localhost:8080/api/users?limit=10&offset=0"
```

#### Create a product
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{
    "name": "iPhone 15",
    "description": "Latest iPhone model",
    "sku": "IPHONE15-128GB",
    "price": 999.99,
    "stock": 50,
    "category_id": "category-uuid",
    "is_active": true
  }'
```

#### OIDC Authentication Flow
```bash
# Step 1: Get authorization URL
curl http://localhost:8080/auth/oidc/login

# Step 2: After user authorization, handle callback
curl "http://localhost:8080/auth/oidc/callback?code=auth-code&state=state-value"

# Step 3: Use token for customer operations
curl -H "Authorization: Bearer <oidc-token>" \
  http://localhost:8080/api/customer/profile
```

### GraphQL Examples

#### Query products with category information
```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { products(limit: 10) { id name price category { name } } }"
  }'
```

#### Create an order
```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation CreateOrder($input: CreateOrderInput!) { createOrder(input: $input) { id orderNumber totalAmount } }",
    "variables": {
      "input": {
        "customerId": "customer-uuid",
        "shippingAddress": "123 Main St",
        "billingAddress": "123 Main St",
        "orderItems": [{"productId": "product-uuid", "quantity": 2}]
      }
    }
  }'
```

### Notification Examples

#### Send Email Notification
```bash
curl -X POST http://localhost:8080/api/notifications/email \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "to": "customer@example.com",
    "subject": "Order Confirmation",
    "body": "Thank you for your order!",
    "html_body": "<h1>Order Confirmed</h1><p>Thank you for your order!</p>"
  }'
```

#### Send SMS Notification
```bash
curl -X POST http://localhost:8080/api/notifications/sms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "phone_number": "+1234567890",
    "message": "Your order has been confirmed!"
  }'
```

#### Send Bulk Email
```bash
curl -X POST http://localhost:8080/api/notifications/bulk/email \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "recipients": ["customer1@example.com", "customer2@example.com"],
    "subject": "Newsletter",
    "body": "Check out our latest products!",
    "html_body": "<h2>Newsletter</h2><p>Check out our latest products!</p>"
  }'
```

## Testing

### Running Tests

The project includes comprehensive unit tests, integration tests, and end-to-end tests.

#### Unit Tests
```bash
# Run all unit tests
make test-unit

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/core/services/...
```

#### E2E Tests
```bash
# Run E2E tests (requires running server)
make test-e2e

# Run all tests (unit + e2e)
make test-all

# Run tests with full validation (linting, security, coverage)
make test-full
```

#### Test Coverage
The project maintains high test coverage with a minimum threshold of 80%. Coverage reports are generated in HTML format.

```bash
# Generate coverage report
make test-coverage
# Opens coverage.html in browser
```

### Test Structure
```
tests/
├── e2e/                    # End-to-end tests
│   └── api_test.go        # API endpoint tests
internal/
├── testutils/             # Test utilities and mocks
│   ├── test_db.go        # Test database setup
│   └── mock_services.go  # Mock service implementations
└── adapters/
    └── notifications/
        ├── email_test.go  # Email adapter tests
        └── sms_test.go    # SMS adapter tests
```

## Deployment

### Docker Deployment

#### Build and Run
```bash
# Build Docker image
make docker-build

# Run with Docker Compose
make docker-up

# Stop services
make docker-down
```

#### Docker Compose Services
- **Application**: Main API server
- **PostgreSQL**: Database
- **Redis**: Caching layer
- **Prometheus**: Metrics collection

### Kubernetes Deployment (Minikube)

#### Prerequisites
- Minikube installed and running
- kubectl configured

#### Deploy to Minikube
```bash
# Deploy entire application stack
make deploy-minikube

# Delete deployment
make delete-minikube
```

#### Kubernetes Resources
- **Namespace**: `sil-backend-assessment`
- **ConfigMap**: Application configuration
- **Secrets**: Sensitive configuration (JWT, SMTP, etc.)
- **Deployments**: Application, PostgreSQL, Redis
- **Services**: Internal and external service exposure
- **Ingress**: HTTP routing with NGINX
- **HPA**: Horizontal Pod Autoscaler

#### Access Points
After deployment, the application is accessible via:
- **NodePort**: `http://<minikube-ip>:30080`
- **Health Check**: `http://<minikube-ip>:30080/api/health`
- **API Docs**: `http://<minikube-ip>:30080/docs`
- **GraphQL Playground**: `http://<minikube-ip>:30080/graphql/playground`

### Production Deployment

#### Environment Configuration
Update the Kubernetes ConfigMap and Secrets with production values:

```bash
# Update secrets
kubectl create secret generic sil-backend-secrets \
  --from-literal=jwt-secret="your-production-jwt-secret" \
  --from-literal=smtp-password="your-production-smtp-password" \
  --from-literal=at-api-key="your-production-at-api-key" \
  --namespace=sil-backend-assessment
```

#### Scaling
The deployment includes Horizontal Pod Autoscaler (HPA) for automatic scaling based on CPU and memory usage.

## Development

### Code Quality

#### Linting and Formatting
```bash
# Run linting
make lint-full

# Format code
make format

# Security scan
make security-scan
```

#### Pre-commit Hooks
The project includes pre-commit hooks for code quality:
- Go formatting
- Linting
- Security scanning
- Test execution

### Adding New Features

1. **Add Domain Models**: Create new entities in `internal/core/domain/`
2. **Define Ports**: Add interfaces in `internal/core/ports/`
3. **Implement Services**: Add business logic in `internal/core/services/`
4. **Create Adapters**: Implement external interfaces in `internal/adapters/`
5. **Add Tests**: Create unit and integration tests
6. **Update Documentation**: Update API documentation and README

### Available Make Commands

```bash
# Development
make build              # Build the application
make run                # Run the application
make deps               # Install dependencies
make clean              # Clean build artifacts

# Testing
make test               # Run basic tests
make test-unit          # Run unit tests only
make test-e2e           # Run E2E tests only
make test-all           # Run all tests
make test-coverage      # Run tests with coverage
make test-full          # Run tests with linting and security

# Code Quality
make lint               # Run basic linting
make lint-full          # Run comprehensive linting
make format             # Format code
make security-scan      # Run security scan

# Database
make migrate_up         # Run database migrations
make migrate_down       # Rollback migrations
make migrate_reset      # Reset database

# Docker
make docker-build       # Build Docker image
make docker-up          # Start with Docker Compose
make docker-down        # Stop Docker Compose

# Kubernetes
make deploy-minikube    # Deploy to Minikube
make delete-minikube    # Delete from Minikube

# GraphQL
make generate-graphql   # Generate GraphQL code
```

## Project Status

### ✅ Completed Features

- **REST API**: Complete CRUD operations for all entities
- **GraphQL API**: Full GraphQL implementation with 37 operations
- **Authentication**: JWT and OIDC authentication
- **Notifications**: Email (SMTP) and SMS (Africa's Talking) support
- **Database**: PostgreSQL with migrations
- **Caching**: Redis integration
- **Testing**: Comprehensive unit and E2E tests
- **Documentation**: Complete API documentation
- **Deployment**: Docker and Kubernetes (Minikube) support
- **Monitoring**: Health checks and metrics

### 📊 API Statistics

- **Total Endpoints**: 45+
- **REST Endpoints**: 30+
- **GraphQL Operations**: 37 (21 queries + 16 mutations)
- **Authentication Methods**: JWT + OIDC
- **Notification Channels**: Email + SMS
- **Test Coverage**: 80%+ threshold

### 🔧 Technical Specifications

- **Language**: Go 1.24+
- **Architecture**: Hexagonal (Ports & Adapters)
- **Database**: PostgreSQL 12+
- **Cache**: Redis 6+
- **HTTP Router**: BunRouter
- **ORM**: Bun
- **GraphQL**: gqlgen
- **Container**: Docker
- **Orchestration**: Kubernetes
- **Testing**: Go testing + custom E2E framework

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License.