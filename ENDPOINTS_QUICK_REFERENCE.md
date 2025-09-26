# API Endpoints Quick Reference

This document provides a quick reference for all available endpoints in the SIL Backend Assessment API.

## Base URLs

- **REST API**: `http://localhost:8080/api`
- **GraphQL**: `http://localhost:8080/graphql`
- **GraphQL Playground**: `http://localhost:8080/graphql/playground`

## Authentication

- **JWT Token**: `Authorization: Bearer <jwt-token>` (for users)
- **OIDC Token**: `Authorization: Bearer <oidc-token>` (for customers)

### Auth Scopes

- ANY: Either a valid User JWT or an OIDC Customer token
- USER: Requires a valid User JWT
- CUSTOMER: Requires a valid OIDC Customer token

## REST API Endpoints

### System
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/` | API information | No |
| GET | `/api/health` | Health check | No |
| GET | `/docs` | API documentation links | No |

### User Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/users` | Create user | JWT |
| GET | `/api/users` | List users (paginated) | JWT |
| GET | `/api/users/{id}` | Get user by ID | JWT |
| PUT | `/api/users/{id}` | Update user | JWT |
| DELETE | `/api/users/{id}` | Delete user | JWT |

### Category Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/categories` | Create category | JWT |
| GET | `/api/categories` | List categories (paginated) | No |
| GET | `/api/categories/{id}` | Get category by ID | No |
| PUT | `/api/categories/{id}` | Update category | JWT |
| DELETE | `/api/categories/{id}` | Delete category | JWT |

### Product Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/products` | Create product | JWT |
| GET | `/api/products` | List products (paginated, filtered) | No |
| GET | `/api/products/{id}` | Get product by ID | No |
| PUT | `/api/products/{id}` | Update product | JWT |
| DELETE | `/api/products/{id}` | Delete product | JWT |

**Query Parameters for GET /api/products:**
- `limit`: Number of products (default: 10)
- `offset`: Skip products (default: 0)
- `category_id`: Filter by category UUID
- `is_active`: Filter by active status (true/false)

### Order Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/orders` | Create order | JWT |
| GET | `/api/orders` | List orders (paginated, filtered) | JWT |
| GET | `/api/orders/{id}` | Get order by ID | JWT |
| PUT | `/api/orders/{id}` | Update order | JWT |
| DELETE | `/api/orders/{id}` | Delete order | JWT |

**Query Parameters for GET /api/orders:**
- `limit`: Number of orders (default: 10)
- `offset`: Skip orders (default: 0)
- `customer_id`: Filter by customer UUID
- `status`: Filter by order status (PENDING, CONFIRMED, PROCESSING, SHIPPED, DELIVERED, CANCELLED)

### Notification Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/notifications/email` | Send email notification | ANY |
| POST | `/api/notifications/sms` | Send SMS notification | ANY |
| POST | `/api/notifications/bulk/email` | Send bulk email notifications | ANY |
| POST | `/api/notifications/bulk/sms` | Send bulk SMS notifications | ANY |
| POST | `/api/notifications/send` | Send generic notification (email/sms) | ANY |

**Notification Features:**
- Email notifications via SMTP (Gmail, Outlook, etc.)
- SMS notifications via Africa's Talking
- Bulk messaging capabilities
- HTML email support
- Phone number validation
- Email address validation

### OIDC Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/auth/oidc/login` | Get authorization URL | No |
| GET | `/auth/oidc/callback` | Handle OIDC callback | No |
| GET | `/auth/oidc/validate` | Validate OIDC token | OIDC Token |
| POST | `/auth/oidc/logout` | OIDC logout | OIDC Token |

### Customer Profile Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/customer/profile` | Get customer profile | OIDC Token |
| PUT | `/api/customer/profile` | Update customer profile | OIDC Token |
| DELETE | `/api/customer/account` | Delete customer account | OIDC Token |

## GraphQL API

### Endpoint
- **URL**: `http://localhost:8080/graphql`
- **Method**: `POST`
- **Content-Type**: `application/json`

### Available Operations

#### Queries
| Operation | Description | Parameters | Scope |
|-----------|-------------|------------|-------|
| `users` | Get users | `pagination` | USER |
| `user` | Get user by ID | `id!` | USER |
| `customers` | Get customers | `pagination` | ANY |
| `customer` | Get customer by ID | `id!` | ANY |
| `categories` | Get categories | `pagination` | ANY |
| `category` | Get category by ID | `id!` | ANY |
| `products` | Get products | `filter`, `pagination` | ANY |
| `product` | Get product by ID | `id!` | ANY |
| `orders` | Get orders | `filter`, `pagination` | ANY |
| `order` | Get order by ID | `id!` | ANY |
| `ordersByCustomer` | Orders by customer | `customerId`, `pagination` | ANY |
| `ordersByStatus` | Orders by status | `status`, `pagination` | USER |
| `orderByNumber` | Order by number | `orderNumber` | ANY |
| `orderStats` | Order statistics | – | USER |
| `productStats` | Product statistics | – | USER |
| `customerStats` | Customer statistics | – | USER |

#### Mutations
| Operation | Description | Input/Args | Scope |
|-----------|-------------|------------|-------|
| `createUser` | Create user | `CreateUserInput!` | USER |
| `updateUser` | Update user | `id!, UpdateUserInput!` | USER |
| `deleteUser` | Delete user | `id!` | USER |
| `createCustomer` | Create customer | `CreateCustomerInput!` | ANY |
| `updateCustomer` | Update customer | `id!, UpdateCustomerInput!` | ANY |
| `deleteCustomer` | Delete customer | `id!` | USER |
| `createCategory` | Create category | `CreateCategoryInput!` | USER |
| `updateCategory` | Update category | `id!, UpdateCategoryInput!` | USER |
| `deleteCategory` | Delete category | `id!` | USER |
| `createProduct` | Create product | `CreateProductInput!` | USER |
| `updateProduct` | Update product | `id!, UpdateProductInput!` | USER |
| `deleteProduct` | Delete product | `id!` | USER |
| `updateProductStock` | Update stock | `id!, stock: Int!` | USER |
| `createOrder` | Create order | `CreateOrderInput!` | ANY |
| `updateOrder` | Update order | `id!, UpdateOrderInput!` | ANY |
| `deleteOrder` | Delete order | `id!` | USER |
| `cancelOrder` | Cancel order | `id!` | ANY |
| `shipOrder` | Mark shipped | `id!` | USER |
| `deliverOrder` | Mark delivered | `id!` | USER |

## Common HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Successful request |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request parameters |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Access denied |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server-side error |

## Quick Examples

### REST API
```bash
# Health check
curl http://localhost:8080/api/health

# Get products with filtering
curl "http://localhost:8080/api/products?limit=5&is_active=true"

# Create a category (requires JWT)
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer <jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Electronics", "description": "Electronic devices"}'

# OIDC login flow
curl http://localhost:8080/auth/oidc/login
```

### GraphQL
```bash
# Query products with category info
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ products(limit: 5) { id name price category { name } } }"
  }'

# Create a user
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { createUser(input: {name: \"John Doe\", email: \"john@example.com\"}) { id name email } }"
  }'
```

## Additional Resources

- **Full API Documentation**: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
- **OIDC Setup Guide**: [OIDC_SETUP.md](./OIDC_SETUP.md)
- **Project Overview**: [README.md](./README.md)
- **GraphQL Playground**: `http://localhost:8080/graphql/playground` (when server is running)