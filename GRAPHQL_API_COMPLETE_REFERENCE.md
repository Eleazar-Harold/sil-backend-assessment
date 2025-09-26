# GraphQL API Complete Reference

## Overview

This document provides a comprehensive reference for all GraphQL operations available in the SIL Backend Assessment API. All operations are protected by JWT/OIDC authentication using the `@auth` directive.

## Authentication

Send your access token using the `Authorization: Bearer <token>` header.

### Auth Scopes
- `@auth(scope: ANY)`: Either a valid User JWT or OIDC Customer token
- `@auth(scope: USER)`: Requires a valid User JWT
- `@auth(scope: CUSTOMER)`: Requires a valid OIDC Customer token

### Endpoints
- **GraphQL Endpoint**: `POST http://localhost:8080/graphql`
- **GraphQL Playground**: `GET http://localhost:8080/graphql/playground`

## Complete Operations List

### Queries (Read Operations)

#### User Queries (USER scope)
| Operation | Description | Parameters | Return Type |
|-----------|-------------|------------|-------------|
| `users` | Get all users with pagination | `pagination: PaginationInput` | `[User!]!` |
| `user` | Get user by ID | `id: ID!` | `User` |
| `searchUsers` | Search users by name/email | `query: String!, pagination: PaginationInput` | `[User!]!` |

#### Customer Queries (ANY scope)
| Operation | Description | Parameters | Return Type |
|-----------|-------------|------------|-------------|
| `customers` | Get all customers with pagination | `pagination: PaginationInput` | `[Customer!]!` |
| `customer` | Get customer by ID | `id: ID!` | `Customer` |
| `searchCustomers` | Search customers by name/email | `query: String!, pagination: PaginationInput` | `[Customer!]!` |

#### Category Queries (ANY scope)
| Operation | Description | Parameters | Return Type |
|-----------|-------------|------------|-------------|
| `categories` | Get all categories with pagination | `pagination: PaginationInput` | `[Category!]!` |
| `category` | Get category by ID | `id: ID!` | `Category` |
| `rootCategories` | Get root categories (no parent) | `pagination: PaginationInput` | `[Category!]!` |
| `subcategories` | Get subcategories of parent | `parentId: ID!, pagination: PaginationInput` | `[Category!]!` |

#### Product Queries (ANY scope)
| Operation | Description | Parameters | Return Type |
|-----------|-------------|------------|-------------|
| `products` | Get products with filtering/pagination | `filter: ProductFilterInput, pagination: PaginationInput` | `[Product!]!` |
| `product` | Get product by ID | `id: ID!` | `Product` |
| `productsByCategory` | Get products by category | `categoryId: ID!, pagination: PaginationInput` | `[Product!]!` |
| `activeProducts` | Get active products only | `pagination: PaginationInput` | `[Product!]!` |
| `searchProducts` | Search products by name/description/SKU | `query: String!, pagination: PaginationInput` | `[Product!]!` |

#### Order Queries
| Operation | Description | Parameters | Return Type | Scope |
|-----------|-------------|------------|-------------|-------|
| `orders` | Get orders with filtering/pagination | `filter: OrderFilterInput, pagination: PaginationInput` | `[Order!]!` | ANY |
| `order` | Get order by ID | `id: ID!` | `Order` | ANY |
| `ordersByCustomer` | Get orders by customer | `customerId: ID!, pagination: PaginationInput` | `[Order!]!` | ANY |
| `ordersByStatus` | Get orders by status | `status: OrderStatus!, pagination: PaginationInput` | `[Order!]!` | USER |
| `orderByNumber` | Get order by order number | `orderNumber: String!` | `Order` | ANY |

#### Statistics Queries (USER scope)
| Operation | Description | Parameters | Return Type |
|-----------|-------------|------------|-------------|
| `orderStats` | Get order statistics | None | `OrderStats!` |
| `productStats` | Get product statistics | None | `ProductStats!` |
| `customerStats` | Get customer statistics | None | `CustomerStats!` |

### Mutations (Write Operations)

#### User Mutations (USER scope)
| Operation | Description | Parameters | Return Type |
|-----------|-------------|------------|-------------|
| `createUser` | Create a new user | `input: CreateUserInput!` | `User!` |
| `updateUser` | Update existing user | `id: ID!, input: UpdateUserInput!` | `User!` |
| `deleteUser` | Delete a user | `id: ID!` | `Boolean!` |

#### Customer Mutations
| Operation | Description | Parameters | Return Type | Scope |
|-----------|-------------|------------|-------------|-------|
| `createCustomer` | Create a new customer | `input: CreateCustomerInput!` | `Customer!` | ANY |
| `updateCustomer` | Update existing customer | `id: ID!, input: UpdateCustomerInput!` | `Customer!` | ANY |
| `deleteCustomer` | Delete a customer | `id: ID!` | `Boolean!` | USER |

#### Category Mutations (USER scope)
| Operation | Description | Parameters | Return Type |
|-----------|-------------|------------|-------------|
| `createCategory` | Create a new category | `input: CreateCategoryInput!` | `Category!` |
| `updateCategory` | Update existing category | `id: ID!, input: UpdateCategoryInput!` | `Category!` |
| `deleteCategory` | Delete a category | `id: ID!` | `Boolean!` |

#### Product Mutations (USER scope)
| Operation | Description | Parameters | Return Type |
|-----------|-------------|------------|-------------|
| `createProduct` | Create a new product | `input: CreateProductInput!` | `Product!` |
| `updateProduct` | Update existing product | `id: ID!, input: UpdateProductInput!` | `Product!` |
| `deleteProduct` | Delete a product | `id: ID!` | `Boolean!` |
| `updateProductStock` | Update product stock | `id: ID!, stock: Int!` | `Product!` |

#### Order Mutations
| Operation | Description | Parameters | Return Type | Scope |
|-----------|-------------|------------|-------------|-------|
| `createOrder` | Create a new order | `input: CreateOrderInput!` | `Order!` | ANY |
| `updateOrder` | Update existing order | `id: ID!, input: UpdateOrderInput!` | `Order!` | ANY |
| `deleteOrder` | Delete an order | `id: ID!` | `Boolean!` | USER |
| `cancelOrder` | Cancel an order | `id: ID!` | `Order!` | ANY |
| `shipOrder` | Mark order as shipped | `id: ID!` | `Order!` | USER |
| `deliverOrder` | Mark order as delivered | `id: ID!` | `Order!` | USER |

## Input Types

### PaginationInput
```graphql
input PaginationInput {
  limit: Int    # Max items (default: 10, max: 100)
  offset: Int   # Items to skip (default: 0)
}
```

### ProductFilterInput
```graphql
input ProductFilterInput {
  categoryId: ID      # Filter by category
  isActive: Boolean   # Filter by active status
  search: String      # Search in name/description
  minPrice: Float     # Minimum price
  maxPrice: Float     # Maximum price
  minStock: Int       # Minimum stock level
}
```

### OrderFilterInput
```graphql
input OrderFilterInput {
  customerId: ID        # Filter by customer
  status: OrderStatus   # Filter by status
  startDate: Time       # Date range start
  endDate: Time         # Date range end
}
```

### Create/Update Inputs

#### User Inputs
```graphql
input CreateUserInput {
  name: String!
  email: String!
}

input UpdateUserInput {
  name: String
  email: String
}
```

#### Customer Inputs
```graphql
input CreateCustomerInput {
  firstName: String!
  lastName: String!
  email: String!
  phone: String
  address: String
  city: String
  state: String
  zipCode: String
  country: String
}

input UpdateCustomerInput {
  firstName: String
  lastName: String
  phone: String
  address: String
  city: String
  state: String
  zipCode: String
  country: String
}
```

#### Category Inputs
```graphql
input CreateCategoryInput {
  name: String!
  description: String
  parentId: ID
}

input UpdateCategoryInput {
  name: String
  description: String
  parentId: ID
}
```

#### Product Inputs
```graphql
input CreateProductInput {
  name: String!
  description: String
  sku: String!
  price: Float!
  stock: Int!
  categoryId: ID!
  isActive: Boolean
}

input UpdateProductInput {
  name: String
  description: String
  sku: String
  price: Float
  stock: Int
  categoryId: ID
  isActive: Boolean
}
```

#### Order Inputs
```graphql
input CreateOrderItemInput {
  productId: ID!
  quantity: Int!
}

input CreateOrderInput {
  customerId: ID!
  shippingAddress: String!
  billingAddress: String!
  notes: String
  orderItems: [CreateOrderItemInput!]!
}

input UpdateOrderInput {
  status: OrderStatus
  shippingAddress: String
  billingAddress: String
  notes: String
  shippedDate: Time
  deliveredDate: Time
}
```

## Return Types

### Core Types
```graphql
type User {
  id: ID!
  name: String!
  email: String!
  createdAt: Time!
  updatedAt: Time!
}

type Customer {
  id: ID!
  firstName: String!
  lastName: String!
  email: String!
  phone: String
  address: String
  city: String
  state: String
  zipCode: String
  country: String
  orders: [Order!]!
  createdAt: Time!
  updatedAt: Time!
}

type Category {
  id: ID!
  name: String!
  description: String
  parentId: ID
  parent: Category
  children: [Category!]!
  products: [Product!]!
  createdAt: Time!
  updatedAt: Time!
}

type Product {
  id: ID!
  name: String!
  description: String
  sku: String!
  price: Float!
  stock: Int!
  categoryId: ID!
  category: Category!
  isActive: Boolean!
  orderItems: [OrderItem!]!
  createdAt: Time!
  updatedAt: Time!
}

type Order {
  id: ID!
  customerId: ID!
  customer: Customer!
  orderNumber: String!
  status: OrderStatus!
  totalAmount: Float!
  shippingAddress: String!
  billingAddress: String!
  notes: String
  orderDate: Time!
  shippedDate: Time
  deliveredDate: Time
  orderItems: [OrderItem!]!
  createdAt: Time!
  updatedAt: Time!
}

type OrderItem {
  id: ID!
  orderId: ID!
  productId: ID!
  product: Product!
  quantity: Int!
  unitPrice: Float!
  totalPrice: Float!
  createdAt: Time!
  updatedAt: Time!
}
```

### Statistics Types
```graphql
type OrderStats {
  totalOrders: Int!
  totalRevenue: Float!
  ordersByStatus: [OrderStatusCount!]!
  averageOrderValue: Float!
  ordersToday: Int!
  revenueToday: Float!
}

type OrderStatusCount {
  status: OrderStatus!
  count: Int!
}

type ProductStats {
  totalProducts: Int!
  activeProducts: Int!
  inactiveProducts: Int!
  lowStockProducts: Int!
  outOfStockProducts: Int!
  totalInventoryValue: Float!
}

type CustomerStats {
  totalCustomers: Int!
  newCustomersThisMonth: Int!
  customersWithOrders: Int!
  topCustomers: [CustomerOrderSummary!]!
}

type CustomerOrderSummary {
  customer: Customer!
  totalOrders: Int!
  totalSpent: Float!
  lastOrderDate: Time
}
```

### Enums
```graphql
enum OrderStatus {
  PENDING
  CONFIRMED
  PROCESSING
  SHIPPED
  DELIVERED
  CANCELLED
}

enum AuthScope {
  ANY
  USER
  CUSTOMER
}
```

## Example Queries

### Basic Queries
```graphql
# Get all products with pagination
query GetProducts {
  products(pagination: { limit: 10, offset: 0 }) {
    id
    name
    price
    stock
    category {
      name
    }
  }
}

# Get orders by customer
query GetCustomerOrders($customerId: ID!) {
  ordersByCustomer(customerId: $customerId, pagination: { limit: 5 }) {
    id
    orderNumber
    status
    totalAmount
    orderItems {
      quantity
      product {
        name
        price
      }
    }
  }
}

# Search products
query SearchProducts($query: String!) {
  searchProducts(query: $query, pagination: { limit: 20 }) {
    id
    name
    description
    price
    sku
  }
}
```

### Statistics Queries
```graphql
query GetDashboardStats {
  orderStats {
    totalOrders
    totalRevenue
    ordersToday
    revenueToday
    ordersByStatus {
      status
      count
    }
  }
  
  productStats {
    totalProducts
    activeProducts
    lowStockProducts
    outOfStockProducts
  }
  
  customerStats {
    totalCustomers
    newCustomersThisMonth
    customersWithOrders
  }
}
```

## Example Mutations

### Creating Entities
```graphql
# Create a customer
mutation CreateCustomer($input: CreateCustomerInput!) {
  createCustomer(input: $input) {
    id
    firstName
    lastName
    email
    createdAt
  }
}

# Create a product
mutation CreateProduct($input: CreateProductInput!) {
  createProduct(input: $input) {
    id
    name
    sku
    price
    stock
    category {
      name
    }
  }
}

# Create an order
mutation CreateOrder($input: CreateOrderInput!) {
  createOrder(input: $input) {
    id
    orderNumber
    status
    totalAmount
    customer {
      firstName
      lastName
    }
    orderItems {
      quantity
      unitPrice
      product {
        name
        sku
      }
    }
  }
}
```

### Updating Entities
```graphql
# Update product stock
mutation UpdateStock($id: ID!, $stock: Int!) {
  updateProductStock(id: $id, stock: $stock) {
    id
    name
    stock
    updatedAt
  }
}

# Ship an order
mutation ShipOrder($id: ID!) {
  shipOrder(id: $id) {
    id
    orderNumber
    status
    shippedDate
  }
}

# Cancel an order
mutation CancelOrder($id: ID!) {
  cancelOrder(id: $id) {
    id
    orderNumber
    status
    updatedAt
  }
}
```

## Testing with cURL

### With OIDC Token
```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_OIDC_TOKEN" \
  -d '{
    "query": "{ products(pagination: { limit: 5 }) { id name price } }"
  }'
```

### With JWT Token
```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "query": "{ users(pagination: { limit: 10 }) { id name email } }"
  }'
```

## GraphQL Playground Usage

1. Open `http://localhost:8080/graphql/playground`
2. Click on the HTTP Headers pane (bottom left)
3. Add your authorization header:
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE"
}
```
4. Write and execute queries/mutations in the main panel

## Error Handling

GraphQL returns errors in the standard format:
```json
{
  "errors": [
    {
      "message": "User not found",
      "locations": [{"line": 2, "column": 3}],
      "path": ["user"]
    }
  ],
  "data": {
    "user": null
  }
}
```

Common error scenarios:
- **401 Unauthorized**: Missing or invalid token
- **403 Forbidden**: Insufficient scope for operation
- **400 Bad Request**: Invalid query syntax or parameters
- **404 Not Found**: Resource doesn't exist
- **500 Internal Server Error**: Server-side error

## Performance Tips

1. **Use specific field selection** - Only request fields you need
2. **Implement pagination** - Always paginate large datasets
3. **Use filters** - Filter data at the query level
4. **Avoid deep nesting** - Limit relationship depth
5. **Cache results** - Implement client-side caching for frequently accessed data

## Complete Operation Count

- **Queries**: 21 total operations
- **Mutations**: 16 total operations
- **Total Operations**: 37

This represents a complete CRUD API for all entities (users, customers, categories, products, orders) with advanced features like search, filtering, pagination, and analytics.
