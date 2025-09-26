# GraphQL API Documentation

## Authentication

All GraphQL operations are protected by JWT/OIDC via the `@auth` directive:

- `@auth(scope: ANY)`: Either a valid User JWT or OIDC Customer token
- `@auth(scope: USER)`: Requires a valid User JWT
- `@auth(scope: CUSTOMER)`: Requires a valid OIDC Customer token

Send your access token using the `Authorization: Bearer <token>` header.

Example curl:

```
curl -sS -H 'Content-Type: application/json' -H 'Authorization: Bearer <TOKEN>' \
  -d '{"query":"{ products(pagination:{limit:5}) { id name } }"}' \
  http://localhost:8080/graphql
```

Playground is available at `/graphql/playground` (you can set HTTP Headers in the playground to include the `Authorization` header).

## Overview

The SIL Backend Assessment provides a comprehensive GraphQL API built with [gqlgen](https://gqlgen.com/) following GraphQL best practices. This API offers a flexible, type-safe, and efficient way to interact with all system entities.

## 🚀 Quick Start

### Endpoints
- **GraphQL Endpoint**: `POST http://localhost:8080/graphql`
- **GraphQL Playground**: `GET http://localhost:8080/graphql/playground`

### Authentication
The GraphQL API supports the same authentication methods as the REST API:
- **JWT Authentication**: For user operations (add `Authorization: Bearer <token>` header)
- **OIDC Authentication**: For customer operations

## 📊 Schema Overview
### Operation Scopes Summary

Use the `Authorization: Bearer <token>` header. Scopes are enforced via `@auth(scope: AuthScope)`:

| Operation | Scope |
|---|---|
| users, user, searchUsers | USER |
| customers, customer | ANY |
| categories, category, rootCategories, subcategories | ANY |
| products, product, productsByCategory, activeProducts, searchProducts | ANY |
| orders, order, ordersByCustomer, orderByNumber | ANY |
| ordersByStatus | USER |
| orderStats, productStats, customerStats | USER |
| create/update/deleteUser | USER |
| create/update/deleteCategory | USER |
| create/update/deleteProduct, updateProductStock | USER |
| create/update/cancelOrder | ANY |
| delete/ship/deliverOrder | USER |


### Core Types

#### User
```graphql
type User {
  id: ID!
  name: String!
  email: String!
  createdAt: Time!
  updatedAt: Time!
}
```

#### Customer
```graphql
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
```

#### Category
```graphql
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
```

#### Product
```graphql
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
```

#### Order
```graphql
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
```

#### OrderItem
```graphql
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

### Enums

#### OrderStatus
```graphql
enum OrderStatus {
  PENDING
  CONFIRMED
  PROCESSING
  SHIPPED
  DELIVERED
  CANCELLED
}
```

## 🔍 Queries

### User Queries (USER scope)

#### Get All Users
```graphql
query GetUsers($pagination: PaginationInput) {
  users(pagination: $pagination) {
    id
    name
    email
    createdAt
    updatedAt
  }
}
```

#### Get User by ID
```graphql
query GetUser($id: ID!) {
  user(id: $id) {
    id
    name
    email
    createdAt
    updatedAt
  }
}
```

#### Search Users
```graphql
query SearchUsers($query: String!, $pagination: PaginationInput) {
  searchUsers(query: $query, pagination: $pagination) {
    id
    name
    email
  }
}
```

### Customer Queries (ANY scope)

#### Get All Customers
```graphql
query GetCustomers($pagination: PaginationInput) {
  customers(pagination: $pagination) {
    id
    firstName
    lastName
    email
    phone
    address
    city
    state
    zipCode
    country
    createdAt
    updatedAt
  }
}
```

#### Get Customer by ID
```graphql
query GetCustomer($id: ID!) {
  customer(id: $id) {
    id
    firstName
    lastName
    email
    phone
    address
    city
    state
    zipCode
    country
    createdAt
    updatedAt
  }
}
```

#### Search Customers
```graphql
query SearchCustomers($query: String!, $pagination: PaginationInput) {
  searchCustomers(query: $query, pagination: $pagination) {
    id
    firstName
    lastName
    email
  }
}
```

#### Get Customer with Orders
```graphql
query GetCustomerWithOrders($id: ID!) {
  customer(id: $id) {
    id
    firstName
    lastName
    email
    orders {
      id
      orderNumber
      status
      totalAmount
      orderDate
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
}
```

### Category Queries (ANY scope)

#### Get All Categories
```graphql
query GetCategories($pagination: PaginationInput) {
  categories(pagination: $pagination) {
    id
    name
    description
    parentId
    parent {
      id
      name
    }
    children {
      id
      name
    }
    products {
      id
      name
      price
    }
    createdAt
    updatedAt
  }
}
```

#### Get Category by ID
```graphql
query GetCategory($id: ID!) {
  category(id: $id) {
    id
    name
    description
    parentId
    parent {
      id
      name
    }
    children {
      id
      name
    }
    products {
      id
      name
      price
    }
    createdAt
    updatedAt
  }
}
```

#### Get Root Categories
```graphql
query GetRootCategories($pagination: PaginationInput) {
  rootCategories(pagination: $pagination) {
    id
    name
    description
    children {
      id
      name
    }
  }
}
```

#### Get Subcategories
```graphql
query GetSubcategories($parentId: ID!, $pagination: PaginationInput) {
  subcategories(parentId: $parentId, pagination: $pagination) {
    id
    name
    description
    products {
      id
      name
      price
    }
  }
}
```

### Product Queries (ANY scope)

#### Get All Products with Filtering
```graphql
query GetProducts($filter: ProductFilterInput, $pagination: PaginationInput) {
  products(filter: $filter, pagination: $pagination) {
    id
    name
    description
    sku
    price
    stock
    isActive
    categoryId
    category {
      id
      name
      parent {
        name
      }
    }
    createdAt
    updatedAt
  }
}
```

#### Get Product by ID
```graphql
query GetProduct($id: ID!) {
  product(id: $id) {
    id
    name
    description
    sku
    price
    stock
    categoryId
    category {
      id
      name
      description
    }
    isActive
    createdAt
    updatedAt
  }
}
```

**Variables:**
```json
{
  "filter": {
    "categoryId": "category-uuid",
    "isActive": true,
    "minPrice": 10.0,
    "maxPrice": 1000.0,
    "minStock": 1
  },
  "pagination": {
    "limit": 20,
    "offset": 0
  }
}
```

#### Get Products by Category
```graphql
query GetProductsByCategory($categoryId: ID!, $pagination: PaginationInput) {
  productsByCategory(categoryId: $categoryId, pagination: $pagination) {
    id
    name
    price
    stock
    category {
      name
    }
  }
}
```

#### Get Active Products Only
```graphql
query GetActiveProducts($pagination: PaginationInput) {
  activeProducts(pagination: $pagination) {
    id
    name
    price
    stock
    isActive
  }
}
```

### Order Queries

#### Get All Orders with Filtering (ANY scope)
```graphql
query GetOrders($filter: OrderFilterInput, $pagination: PaginationInput) {
  orders(filter: $filter, pagination: $pagination) {
    id
    customerId
    orderNumber
    status
    totalAmount
    shippingAddress
    billingAddress
    notes
    orderDate
    shippedDate
    deliveredDate
    customer {
      firstName
      lastName
      email
    }
    orderItems {
      id
      productId
      quantity
      unitPrice
      totalPrice
      product {
        id
        name
        sku
        price
      }
    }
    createdAt
    updatedAt
  }
}
```

#### Get Order by ID (ANY scope)
```graphql
query GetOrder($id: ID!) {
  order(id: $id) {
    id
    customerId
    orderNumber
    status
    totalAmount
    shippingAddress
    billingAddress
    notes
    orderDate
    shippedDate
    deliveredDate
    customer {
      id
      firstName
      lastName
      email
    }
    orderItems {
      id
      productId
      quantity
      unitPrice
      totalPrice
      product {
        id
        name
        sku
        price
      }
    }
    createdAt
    updatedAt
  }
}
```

**Variables:**
```json
{
  "filter": {
    "customerId": "customer-uuid",
    "status": "PENDING",
    "startDate": "2024-01-01T00:00:00Z",
    "endDate": "2024-12-31T23:59:59Z"
  },
  "pagination": {
    "limit": 10,
    "offset": 0
  }
}
```

#### Get Orders by Customer (ANY scope)
```graphql
query GetOrdersByCustomer($customerId: ID!, $pagination: PaginationInput) {
  ordersByCustomer(customerId: $customerId, pagination: $pagination) {
    id
    customerId
    orderNumber
    status
    totalAmount
    orderDate
    shippedDate
    deliveredDate
    orderItems {
      id
      productId
      quantity
      unitPrice
      totalPrice
      product {
        id
        name
        sku
        price
      }
    }
    createdAt
    updatedAt
  }
}
```

#### Get Orders by Status (USER scope)
```graphql
query GetOrdersByStatus($status: OrderStatus!, $pagination: PaginationInput) {
  ordersByStatus(status: $status, pagination: $pagination) {
    id
    customerId
    orderNumber
    status
    totalAmount
    orderDate
    customer {
      firstName
      lastName
      email
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

#### Get Order by Number (ANY scope)
```graphql
query GetOrderByNumber($orderNumber: String!) {
  orderByNumber(orderNumber: $orderNumber) {
    id
    customerId
    orderNumber
    status
    totalAmount
    shippingAddress
    billingAddress
    notes
    orderDate
    shippedDate
    deliveredDate
    customer {
      id
      firstName
      lastName
      email
    }
    orderItems {
      id
      productId
      quantity
      unitPrice
      totalPrice
      product {
        id
        name
        sku
        price
      }
    }
    createdAt
    updatedAt
  }
}
```

### Statistics Queries (USER scope)

#### Get Order Statistics
```graphql
query GetOrderStats {
  orderStats {
    totalOrders
    totalRevenue
    averageOrderValue
    ordersToday
    revenueToday
    ordersByStatus {
      status
      count
    }
  }
}
```

#### Get Product Statistics
```graphql
query GetProductStats {
  productStats {
    totalProducts
    activeProducts
    inactiveProducts
    lowStockProducts
    outOfStockProducts
    totalInventoryValue
  }
}
```

#### Get Customer Statistics
```graphql
query GetCustomerStats {
  customerStats {
    totalCustomers
    newCustomersThisMonth
    customersWithOrders
    topCustomers {
      customer {
        firstName
        lastName
        email
      }
      totalOrders
      totalSpent
      lastOrderDate
    }
  }
}
```

## ✏️ Mutations

### User Mutations (USER scope)

#### Create User
```graphql
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
    name
    email
    createdAt
  }
}
```

**Variables:**
```json
{
  "input": {
    "name": "John Doe",
    "email": "john.doe@example.com"
  }
}
```

#### Update User
```graphql
mutation UpdateUser($id: ID!, $input: UpdateUserInput!) {
  updateUser(id: $id, input: $input) {
    id
    name
    email
    updatedAt
  }
}
```

#### Delete User
```graphql
mutation DeleteUser($id: ID!) {
  deleteUser(id: $id)
}
```

### Customer Mutations

#### Create Customer (ANY scope)
```graphql
mutation CreateCustomer($input: CreateCustomerInput!) {
  createCustomer(input: $input) {
    id
    firstName
    lastName
    email
    phone
    address
    city
    state
    zipCode
    country
    createdAt
    updatedAt
  }
}
```

#### Update Customer (ANY scope)
```graphql
mutation UpdateCustomer($id: ID!, $input: UpdateCustomerInput!) {
  updateCustomer(id: $id, input: $input) {
    id
    firstName
    lastName
    email
    phone
    address
    city
    state
    zipCode
    country
    createdAt
    updatedAt
  }
}
```

#### Delete Customer (USER scope)
```graphql
mutation DeleteCustomer($id: ID!) {
  deleteCustomer(id: $id)
}
```

**Variables:**
```json
{
  "input": {
    "firstName": "Jane",
    "lastName": "Smith",
    "email": "jane.smith@example.com",
    "phone": "+1-555-0123",
    "address": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zipCode": "10001",
    "country": "USA"
  }
}
```

### Category Mutations (USER scope)

#### Create Category
```graphql
mutation CreateCategory($input: CreateCategoryInput!) {
  createCategory(input: $input) {
    id
    name
    description
    parentId
    createdAt
  }
}
```

**Variables:**
```json
{
  "input": {
    "name": "Electronics",
    "description": "Electronic devices and accessories",
    "parentId": null
  }
}
```

#### Create Subcategory
```graphql
mutation CreateSubcategory($input: CreateCategoryInput!) {
  createCategory(input: $input) {
    id
    name
    description
    parent {
      name
    }
  }
}
```

**Variables:**
```json
{
  "input": {
    "name": "Smartphones",
    "description": "Mobile phones and accessories",
    "parentId": "electronics-category-uuid"
  }
}
```

### Product Mutations (USER scope)

#### Create Product
```graphql
mutation CreateProduct($input: CreateProductInput!) {
  createProduct(input: $input) {
    id
    name
    description
    sku
    price
    stock
    isActive
    category {
      name
    }
    createdAt
  }
}
```

**Variables:**
```json
{
  "input": {
    "name": "iPhone 15 Pro",
    "description": "Latest iPhone with advanced features",
    "sku": "IPHONE15PRO-256GB",
    "price": 1199.99,
    "stock": 50,
    "categoryId": "smartphones-category-uuid",
    "isActive": true
  }
}
```

#### Update Product
```graphql
mutation UpdateProduct($id: ID!, $input: UpdateProductInput!) {
  updateProduct(id: $id, input: $input) {
    id
    name
    description
    sku
    price
    stock
    categoryId
    isActive
    createdAt
    updatedAt
  }
}
```

#### Delete Product
```graphql
mutation DeleteProduct($id: ID!) {
  deleteProduct(id: $id)
}
```

#### Update Product Stock
```graphql
mutation UpdateProductStock($id: ID!, $stock: Int!) {
  updateProductStock(id: $id, stock: $stock) {
    id
    name
    stock
    updatedAt
  }
}
```

### Order Mutations

#### Create Order (ANY scope)
```graphql
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
      totalPrice
      product {
        name
        sku
      }
    }
    createdAt
  }
}
```

**Variables:**
```json
{
  "input": {
    "customerId": "customer-uuid",
    "shippingAddress": "123 Main St, New York, NY 10001",
    "billingAddress": "123 Main St, New York, NY 10001",
    "notes": "Please handle with care",
    "orderItems": [
      {
        "productId": "product-uuid-1",
        "quantity": 2
      },
      {
        "productId": "product-uuid-2",
        "quantity": 1
      }
    ]
  }
}
```

#### Update Order (ANY scope)
```graphql
mutation UpdateOrder($id: ID!, $input: UpdateOrderInput!) {
  updateOrder(id: $id, input: $input) {
    id
    customerId
    orderNumber
    status
    totalAmount
    shippingAddress
    billingAddress
    notes
    orderDate
    shippedDate
    deliveredDate
    createdAt
    updatedAt
  }
}
```

#### Delete Order (USER scope)
```graphql
mutation DeleteOrder($id: ID!) {
  deleteOrder(id: $id)
}
```

#### Ship Order (USER scope)
```graphql
mutation ShipOrder($id: ID!) {
  shipOrder(id: $id) {
    id
    orderNumber
    status
    shippedDate
    updatedAt
  }
}
```

#### Deliver Order (USER scope)
```graphql
mutation DeliverOrder($id: ID!) {
  deliverOrder(id: $id) {
    id
    orderNumber
    status
    deliveredDate
    updatedAt
  }
}
```

#### Cancel Order (ANY scope)
```graphql
mutation CancelOrder($id: ID!) {
  cancelOrder(id: $id) {
    id
    orderNumber
    status
    updatedAt
  }
}
```

## 🔧 Advanced Features

### Pagination

All list queries support pagination with the `PaginationInput` type:

```graphql
input PaginationInput {
  limit: Int    # Maximum items to return (default: 10, max: 100)
  offset: Int   # Number of items to skip (default: 0)
}
```

**Example:**
```graphql
query GetProductsPaginated {
  products(pagination: { limit: 20, offset: 40 }) {
    id
    name
    price
  }
}
```

### Filtering

#### Product Filtering
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

#### Order Filtering
```graphql
input OrderFilterInput {
  customerId: ID        # Filter by customer
  status: OrderStatus   # Filter by status
  startDate: Time       # Date range start
  endDate: Time         # Date range end
}
```

### Relationship Loading

GraphQL automatically resolves relationships. You can query nested data in a single request:

```graphql
query GetCompleteOrderInfo($id: ID!) {
  order(id: $id) {
    id
    orderNumber
    status
    totalAmount
    customer {
      firstName
      lastName
      email
      address
    }
    orderItems {
      quantity
      unitPrice
      totalPrice
      product {
        name
        sku
        price
        category {
          name
          parent {
            name
          }
        }
      }
    }
  }
}
```

## 🛠️ Development Tools

### GraphQL Playground

Access the interactive playground at `http://localhost:8080/graphql/playground` to:
- Explore the schema with auto-completion
- Test queries and mutations
- View real-time documentation
- Debug and optimize queries

### Schema Introspection

GraphQL supports introspection queries to explore the schema programmatically:

```graphql
query IntrospectionQuery {
  __schema {
    types {
      name
      description
      fields {
        name
        type {
          name
        }
      }
    }
  }
}
```

### Error Handling

GraphQL errors follow the standard format:

```json
{
  "errors": [
    {
      "message": "User not found",
      "locations": [
        {
          "line": 2,
          "column": 3
        }
      ],
      "path": ["user"]
    }
  ],
  "data": {
    "user": null
  }
}
```

## 🚀 Performance Optimization

### Query Optimization Tips

1. **Use Specific Field Selection**: Only request fields you need
   ```graphql
   # Good
   query {
     products {
       id
       name
       price
     }
   }
   
   # Avoid
   query {
     products {
       id
       name
       price
       description
       sku
       stock
       category {
         # ... all fields
       }
     }
   }
   ```

2. **Use Pagination**: Always paginate large datasets
   ```graphql
   query {
     products(pagination: { limit: 20, offset: 0 }) {
       id
       name
     }
   }
   ```

3. **Avoid Deep Nesting**: Limit relationship depth
   ```graphql
   # Good
   query {
     categories {
       name
       products {
         name
       }
     }
   }
   
   # Avoid deep nesting
   query {
     categories {
       children {
         children {
           children {
             products {
               category {
                 parent {
                   # ... too deep
                 }
               }
             }
           }
         }
       }
     }
   }
   ```

### Caching

The GraphQL server includes:
- **Query Caching**: Parsed queries are cached for performance
- **Automatic Persisted Queries**: Reduces bandwidth usage
- **Field-level Caching**: Consider implementing for expensive operations

## 🔒 Security Best Practices

### Authentication

Include authentication headers in your requests:

```javascript
// For JWT authentication
fetch('http://localhost:8080/graphql', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-jwt-token'
  },
  body: JSON.stringify({
    query: '{ users { id name } }'
  })
})
```

### Input Validation

All inputs are validated at multiple levels:
- GraphQL schema validation
- Business logic validation
- Database constraints

### Rate Limiting

Consider implementing:
- Query complexity analysis
- Depth limiting
- Request rate limiting

## 📚 Examples

### Complete E-commerce Flow

```graphql
# 1. Create a customer
mutation CreateCustomer {
  createCustomer(input: {
    firstName: "John"
    lastName: "Doe"
    email: "john@example.com"
    address: "123 Main St"
    city: "New York"
    state: "NY"
    zipCode: "10001"
  }) {
    id
    firstName
    lastName
  }
}

# 2. Browse products by category
query BrowseProducts {
  productsByCategory(categoryId: "electronics-uuid") {
    id
    name
    price
    stock
    category {
      name
    }
  }
}

# 3. Create an order
mutation CreateOrder {
  createOrder(input: {
    customerId: "customer-uuid"
    shippingAddress: "123 Main St, New York, NY 10001"
    billingAddress: "123 Main St, New York, NY 10001"
    orderItems: [
      { productId: "product-1-uuid", quantity: 2 }
      { productId: "product-2-uuid", quantity: 1 }
    ]
  }) {
    id
    orderNumber
    totalAmount
    status
  }
}

# 4. Track order status
query TrackOrder {
  orderByNumber(orderNumber: "ORD-2024-001") {
    id
    status
    orderDate
    shippedDate
    deliveredDate
    orderItems {
      quantity
      product {
        name
      }
    }
  }
}
```

### Admin Dashboard Queries

```graphql
# Dashboard overview
query AdminDashboard {
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

# Recent orders
query RecentOrders {
  orders(pagination: { limit: 10 }) {
    id
    orderNumber
    status
    totalAmount
    orderDate
    customer {
      firstName
      lastName
    }
  }
}

# Low stock products
query LowStockProducts {
  products(filter: { maxStock: 10 }) {
    id
    name
    stock
    category {
      name
    }
  }
}
```

## 🤝 Contributing

When extending the GraphQL API:

1. **Update Schema**: Modify `schema.graphqls`
2. **Regenerate Code**: Run `go run github.com/99designs/gqlgen generate`
3. **Implement Resolvers**: Add resolver logic
4. **Update Documentation**: Update this file
5. **Add Tests**: Test new functionality

## 📖 Additional Resources

- [GraphQL Specification](https://spec.graphql.org/)
- [gqlgen Documentation](https://gqlgen.com/)
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)
- [Apollo GraphQL Guide](https://www.apollographql.com/docs/)

## 📊 Complete Operations Summary

### Total Operations: 37
- **Queries**: 21 operations
- **Mutations**: 16 operations

### Operations by Entity
| Entity | Queries | Mutations | Total |
|--------|---------|-----------|-------|
| Users | 3 | 3 | 6 |
| Customers | 3 | 3 | 6 |
| Categories | 4 | 3 | 7 |
| Products | 5 | 4 | 9 |
| Orders | 5 | 6 | 11 |
| Statistics | 3 | 0 | 3 |

### Operations by Scope
| Scope | Queries | Mutations | Total |
|-------|---------|-----------|-------|
| ANY | 12 | 4 | 16 |
| USER | 9 | 12 | 21 |

### Key Features
- ✅ Complete CRUD operations for all entities
- ✅ Advanced filtering and search capabilities
- ✅ Pagination support for all list operations
- ✅ Hierarchical category support
- ✅ Order status management with workflow mutations
- ✅ Comprehensive statistics and analytics
- ✅ JWT/OIDC authentication with scope-based authorization
- ✅ Relationship loading and nested queries
- ✅ Input validation and error handling

### Authentication Summary
- **ANY scope**: Accessible with either User JWT or Customer OIDC token
- **USER scope**: Requires User JWT token only
- **CUSTOMER scope**: Requires Customer OIDC token only (currently not used)

---

For more information about the REST API, see [API_DOCUMENTATION.md](./API_DOCUMENTATION.md).
For quick reference, see [ENDPOINTS_QUICK_REFERENCE.md](./ENDPOINTS_QUICK_REFERENCE.md).
For complete operation reference, see [GRAPHQL_API_COMPLETE_REFERENCE.md](./GRAPHQL_API_COMPLETE_REFERENCE.md).