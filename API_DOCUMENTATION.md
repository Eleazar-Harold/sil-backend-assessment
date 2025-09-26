# SIL Backend Assessment - API Documentation

This document provides comprehensive documentation for all REST and GraphQL endpoints available in the SIL Backend Assessment application.

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Base URLs](#base-urls)
- [REST API Endpoints](#rest-api-endpoints)
- [GraphQL API](#graphql-api)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Overview

The SIL Backend Assessment API provides both REST and GraphQL interfaces for managing users, customers, categories, products, and orders. The API supports two authentication methods:

1. **JWT Authentication** - For traditional user authentication
2. **OIDC Authentication** - For customer authentication via external providers

## Authentication

### JWT Authentication (Users)
Traditional JWT-based authentication for internal users and admin operations.

### OIDC Authentication (Customers)
OpenID Connect authentication for customers using external identity providers (Google, Microsoft, Auth0, etc.).

**Authorization Header Format:**
```
Authorization: Bearer <token>
```

### Auth Scopes

- ANY: Either a valid User JWT or an OIDC Customer token
- USER: Requires a valid User JWT
- CUSTOMER: Requires a valid OIDC Customer token

### REST Endpoints and Required Scopes (summary)

| Endpoint | Method | Description | Scope |
|---|---|---|---|
| /api/users | GET/POST/PUT/DELETE | User CRUD | USER |
| /api/customers | GET/POST/PUT/DELETE | Customer CRUD/auth | ANY (GET/POST/PUT), USER (DELETE) |
| /api/categories | GET/PUT/POST/DELETE | Category CRUD | USER (write), ANY (read) |
| /api/products | GET/PUT/POST/DELETE | Product CRUD | USER (write), ANY (read) |
| /api/orders | GET/PUT/POST/DELETE | Order CRUD | ANY (most), USER (DELETE) |
| /api/notifications/* | POST | Email/SMS notifications | ANY |
| /auth/oidc/* | GET/POST | OIDC auth flow | Public (login/callback), ANY (validate/logout) |

See `docs/ENDPOINTS_QUICK_REFERENCE.md` for more details.

### GraphQL Endpoint and Scopes (summary)

Path: `/graphql`. Playground at `/graphql/playground` (public UI; set Authorization header in the playground).

| Operation (root) | Scope |
|---|---|
| users, user, searchUsers | USER |
| customers, customer, searchCustomers | ANY |
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

Example headers:

```
Authorization: Bearer <token>
Content-Type: application/json
```

Example cURL (GraphQL):

```
curl -sS -H 'Content-Type: application/json' -H 'Authorization: Bearer <TOKEN>' \
  -d '{"query":"{ products(pagination:{limit:5}) { id name } }"}' \
  http://localhost:8080/graphql
```

## Base URLs

- **REST API**: `http://localhost:8080/api`
- **GraphQL API**: `http://localhost:8080/graphql`
- **GraphQL Playground**: `http://localhost:8080/graphql/playground`
- **API Documentation**: `http://localhost:8080/docs`

## REST API Endpoints

### Health Check

#### Get Health Status
- **Endpoint**: `GET /api/health`
- **Description**: Check API health status
- **Authentication**: None required

**Response:**
```json
{
  "status": "ok",
  "service": "sil-backend-assessment"
}
```

### User Management

#### Create User
- **Endpoint**: `POST /api/users`
- **Description**: Create a new user
- **Authentication**: JWT required

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

**Response:**
```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Get User
- **Endpoint**: `GET /api/users/{id}`
- **Description**: Retrieve a specific user by ID
- **Authentication**: JWT required

**Response:**
```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Get Users
- **Endpoint**: `GET /api/users`
- **Description**: Retrieve all users with pagination
- **Authentication**: JWT required
- **Query Parameters**:
  - `limit` (optional): Number of users to return (default: 10)
  - `offset` (optional): Number of users to skip (default: 0)

**Response:**
```json
{
  "users": [
    {
      "id": "uuid",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

#### Update User
- **Endpoint**: `PUT /api/users/{id}`
- **Description**: Update an existing user
- **Authentication**: JWT required

**Request Body:**
```json
{
  "name": "Jane Doe",
  "email": "jane.doe@example.com"
}
```

#### Delete User
- **Endpoint**: `DELETE /api/users/{id}`
- **Description**: Delete a user
- **Authentication**: JWT required

**Response:**
```json
{
  "message": "User deleted successfully"
}
```

### Category Management

#### Create Category
- **Endpoint**: `POST /api/categories`
- **Description**: Create a new category
- **Authentication**: JWT required

**Request Body:**
```json
{
  "name": "Electronics",
  "description": "Electronic devices and accessories",
  "parent_id": "uuid" // optional
}
```

#### Get Category
- **Endpoint**: `GET /api/categories/{id}`
- **Description**: Retrieve a specific category by ID
- **Authentication**: None required

#### Get Categories
- **Endpoint**: `GET /api/categories`
- **Description**: Retrieve all categories with pagination
- **Authentication**: None required
- **Query Parameters**:
  - `limit` (optional): Number of categories to return (default: 10)
  - `offset` (optional): Number of categories to skip (default: 0)

#### Update Category
- **Endpoint**: `PUT /api/categories/{id}`
- **Description**: Update an existing category
- **Authentication**: JWT required

#### Delete Category
- **Endpoint**: `DELETE /api/categories/{id}`
- **Description**: Delete a category
- **Authentication**: JWT required

### Product Management

#### Create Product
- **Endpoint**: `POST /api/products`
- **Description**: Create a new product
- **Authentication**: JWT required

**Request Body:**
```json
{
  "name": "iPhone 15",
  "description": "Latest iPhone model",
  "sku": "IPHONE15-128GB",
  "price": 999.99,
  "stock": 50,
  "category_id": "uuid",
  "is_active": true
}
```

#### Get Product
- **Endpoint**: `GET /api/products/{id}`
- **Description**: Retrieve a specific product by ID
- **Authentication**: None required

#### Get Products
- **Endpoint**: `GET /api/products`
- **Description**: Retrieve all products with pagination and filtering
- **Authentication**: None required
- **Query Parameters**:
  - `limit` (optional): Number of products to return (default: 10)
  - `offset` (optional): Number of products to skip (default: 0)
  - `category_id` (optional): Filter by category ID
  - `is_active` (optional): Filter by active status (true/false)

#### Update Product
- **Endpoint**: `PUT /api/products/{id}`
- **Description**: Update an existing product
- **Authentication**: JWT required

#### Delete Product
- **Endpoint**: `DELETE /api/products/{id}`
- **Description**: Delete a product
- **Authentication**: JWT required

### Order Management

#### Create Order
- **Endpoint**: `POST /api/orders`
- **Description**: Create a new order
- **Authentication**: JWT required

**Request Body:**
```json
{
  "customer_id": "uuid",
  "shipping_address": "123 Main St, City, State 12345",
  "billing_address": "123 Main St, City, State 12345",
  "notes": "Please handle with care",
  "order_items": [
    {
      "product_id": "uuid",
      "quantity": 2
    }
  ]
}
```

#### Get Order
- **Endpoint**: `GET /api/orders/{id}`
- **Description**: Retrieve a specific order by ID
- **Authentication**: JWT required

#### Get Orders
- **Endpoint**: `GET /api/orders`
- **Description**: Retrieve all orders with pagination and filtering
- **Authentication**: JWT required
- **Query Parameters**:
  - `limit` (optional): Number of orders to return (default: 10)
  - `offset` (optional): Number of orders to skip (default: 0)
  - `customer_id` (optional): Filter by customer ID
  - `status` (optional): Filter by order status (PENDING, CONFIRMED, PROCESSING, SHIPPED, DELIVERED, CANCELLED)

#### Update Order
- **Endpoint**: `PUT /api/orders/{id}`
- **Description**: Update an existing order
- **Authentication**: JWT required

**Request Body:**
```json
{
  "status": "SHIPPED",
  "shipping_address": "Updated address",
  "billing_address": "Updated billing address",
  "notes": "Updated notes",
  "shipped_date": "2024-01-01T00:00:00Z",
  "delivered_date": "2024-01-02T00:00:00Z"
}
```

#### Delete Order
- **Endpoint**: `DELETE /api/orders/{id}`
- **Description**: Delete an order
- **Authentication**: JWT required

### OIDC Authentication

#### Get Authorization URL
- **Endpoint**: `GET /auth/oidc/login`
- **Description**: Generate OIDC authorization URL for customer login
- **Authentication**: None required

**Response:**
```json
{
  "auth_url": "https://accounts.google.com/oauth/authorize?...",
  "state": "random-state-string"
}
```

#### Handle OIDC Callback
- **Endpoint**: `GET /auth/oidc/callback`
- **Description**: Handle OIDC provider callback
- **Authentication**: None required
- **Query Parameters**:
  - `code`: Authorization code from OIDC provider
  - `state`: State parameter for CSRF protection

**Response:**
```json
{
  "access_token": "jwt-access-token",
  "refresh_token": "jwt-refresh-token",
  "customer": {
    "id": "uuid",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone": "",
    "address": "",
    "city": "",
    "state": "",
    "zip_code": "",
    "country": "",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "expires_at": "2024-01-01T06:00:00Z",
  "is_new_user": true
}
```

#### Validate OIDC Token
- **Endpoint**: `GET /auth/oidc/validate`
- **Description**: Validate an OIDC token
- **Authentication**: Bearer token required

#### OIDC Logout
- **Endpoint**: `POST /auth/oidc/logout`
- **Description**: Logout from OIDC session
- **Authentication**: Bearer token required

### Customer Profile Management

#### Get Customer Profile
- **Endpoint**: `GET /api/customer/profile`
- **Description**: Get current customer's profile
- **Authentication**: OIDC Bearer token required

**Response:**
```json
{
  "id": "uuid",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "address": "123 Main St",
  "city": "New York",
  "state": "NY",
  "zip_code": "10001",
  "country": "USA",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Update Customer Profile
- **Endpoint**: `PUT /api/customer/profile`
- **Description**: Update current customer's profile
- **Authentication**: OIDC Bearer token required

**Request Body:**
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "phone": "+1987654321",
  "address": "456 Oak Ave",
  "city": "Los Angeles",
  "state": "CA",
  "zip_code": "90210",
  "country": "USA"
}
```

#### Delete Customer Account
- **Endpoint**: `DELETE /api/customer/account`
- **Description**: Delete current customer's account
- **Authentication**: OIDC Bearer token required

**Response:**
```json
{
  "message": "Account deleted successfully"
}
```

### Notification Endpoints

The notification endpoints provide email and SMS capabilities using SMTP and Africa's Talking services.

#### Send Email
- **Endpoint**: `POST /api/notifications/email`
- **Description**: Send an email notification to a single recipient
- **Authentication**: ANY (JWT or OIDC token required)

**Request Body:**
```json
{
  "to": "recipient@example.com",
  "subject": "Email Subject",
  "body": "Plain text email body",
  "html_body": "<p>HTML email body (optional)</p>"
}
```

**Response:**
```json
{
  "message": "Email sent successfully",
  "to": "recipient@example.com"
}
```

#### Send SMS
- **Endpoint**: `POST /api/notifications/sms`
- **Description**: Send an SMS notification to a single phone number
- **Authentication**: ANY (JWT or OIDC token required)

**Request Body:**
```json
{
  "phone_number": "+1234567890",
  "message": "Your SMS message here"
}
```

**Response:**
```json
{
  "message": "SMS sent successfully",
  "phone_number": "+1234567890"
}
```

#### Send Bulk Email
- **Endpoint**: `POST /api/notifications/bulk/email`
- **Description**: Send the same email to multiple recipients
- **Authentication**: ANY (JWT or OIDC token required)

**Request Body:**
```json
{
  "recipients": ["customer1@example.com", "customer2@example.com"],
  "subject": "Newsletter Subject",
  "body": "Newsletter content",
  "html_body": "<h1>Newsletter</h1><p>Content here</p>"
}
```

**Response:**
```json
{
  "message": "Bulk emails sent successfully",
  "recipients": ["customer1@example.com", "customer2@example.com"],
  "count": 2
}
```

#### Send Bulk SMS
- **Endpoint**: `POST /api/notifications/bulk/sms`
- **Description**: Send the same SMS to multiple phone numbers
- **Authentication**: ANY (JWT or OIDC token required)

**Request Body:**
```json
{
  "phone_numbers": ["+1234567890", "+0987654321"],
  "message": "Your bulk SMS message"
}
```

**Response:**
```json
{
  "message": "Bulk SMS sent successfully",
  "phone_numbers": ["+1234567890", "+0987654321"],
  "count": 2
}
```

#### Send Generic Notification
- **Endpoint**: `POST /api/notifications/send`
- **Description**: Send a notification based on type (email or SMS)
- **Authentication**: ANY (JWT or OIDC token required)

**Request Body:**
```json
{
  "type": "email",  // or "sms"
  "to": "recipient@example.com",  // for email
  "phone_number": "+1234567890",  // for SMS
  "subject": "Email Subject",     // for email
  "body": "Email body",          // for email
  "html_body": "<p>HTML body</p>", // for email (optional)
  "message": "SMS message"       // for SMS
}
```

**Response:**
```json
{
  "message": "Notification sent successfully",
  "type": "email"
}
```

For detailed notification API documentation, see [NOTIFICATION_API_DOCUMENTATION.md](./NOTIFICATION_API_DOCUMENTATION.md).

## GraphQL API

### Endpoint
- **URL**: `http://localhost:8080/graphql`
- **Method**: `POST`
- **Content-Type**: `application/json`

### Schema Overview

The GraphQL API provides queries and mutations for all entities:

#### Types
- `User`: User entity with id, name, email, timestamps
- `Customer`: Customer entity with personal and address information
- `Category`: Product category with hierarchical support
- `Product`: Product entity with pricing and inventory
- `Order`: Order entity with items and status tracking
- `OrderItem`: Individual items within an order

#### Enums
- `OrderStatus`: PENDING, CONFIRMED, PROCESSING, SHIPPED, DELIVERED, CANCELLED

### Queries

#### User Queries
```graphql
# Get all users with pagination
query GetUsers($limit: Int, $offset: Int) {
  users(limit: $limit, offset: $offset) {
    id
    name
    email
    createdAt
    updatedAt
  }
}

# Get specific user
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

#### Customer Queries
```graphql
# Get all customers with pagination
query GetCustomers($limit: Int, $offset: Int) {
  customers(limit: $limit, offset: $offset) {
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

# Get specific customer
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

#### Category Queries
```graphql
# Get all categories with pagination
query GetCategories($limit: Int, $offset: Int) {
  categories(limit: $limit, offset: $offset) {
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

# Get specific category
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
      stock
    }
    createdAt
    updatedAt
  }
}
```

#### Product Queries
```graphql
# Get all products with filtering
query GetProducts($limit: Int, $offset: Int, $categoryId: ID, $isActive: Boolean) {
  products(limit: $limit, offset: $offset, categoryId: $categoryId, isActive: $isActive) {
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
    }
    isActive
    createdAt
    updatedAt
  }
}

# Get specific product
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

#### Order Queries
```graphql
# Get all orders with filtering
query GetOrders($limit: Int, $offset: Int, $customerId: ID, $status: OrderStatus) {
  orders(limit: $limit, offset: $offset, customerId: $customerId, status: $status) {
    id
    customerId
    customer {
      id
      firstName
      lastName
      email
    }
    orderNumber
    status
    totalAmount
    shippingAddress
    billingAddress
    notes
    orderDate
    shippedDate
    deliveredDate
    orderItems {
      id
      productId
      product {
        id
        name
        price
      }
      quantity
      unitPrice
      totalPrice
    }
    createdAt
    updatedAt
  }
}

# Get specific order
query GetOrder($id: ID!) {
  order(id: $id) {
    id
    customerId
    customer {
      id
      firstName
      lastName
      email
      phone
      address
    }
    orderNumber
    status
    totalAmount
    shippingAddress
    billingAddress
    notes
    orderDate
    shippedDate
    deliveredDate
    orderItems {
      id
      productId
      product {
        id
        name
        description
        sku
        price
      }
      quantity
      unitPrice
      totalPrice
      createdAt
      updatedAt
    }
    createdAt
    updatedAt
  }
}
```

### Mutations

#### User Mutations
```graphql
# Create user
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
    name
    email
    createdAt
    updatedAt
  }
}

# Update user
mutation UpdateUser($id: ID!, $input: UpdateUserInput!) {
  updateUser(id: $id, input: $input) {
    id
    name
    email
    createdAt
    updatedAt
  }
}

# Delete user
mutation DeleteUser($id: ID!) {
  deleteUser(id: $id)
}
```

#### Customer Mutations
```graphql
# Create customer
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

# Update customer
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

# Delete customer
mutation DeleteCustomer($id: ID!) {
  deleteCustomer(id: $id)
}
```

#### Category Mutations
```graphql
# Create category
mutation CreateCategory($input: CreateCategoryInput!) {
  createCategory(input: $input) {
    id
    name
    description
    parentId
    createdAt
    updatedAt
  }
}

# Update category
mutation UpdateCategory($id: ID!, $input: UpdateCategoryInput!) {
  updateCategory(id: $id, input: $input) {
    id
    name
    description
    parentId
    createdAt
    updatedAt
  }
}

# Delete category
mutation DeleteCategory($id: ID!) {
  deleteCategory(id: $id)
}
```

#### Product Mutations
```graphql
# Create product
mutation CreateProduct($input: CreateProductInput!) {
  createProduct(input: $input) {
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
    }
    isActive
    createdAt
    updatedAt
  }
}

# Update product
mutation UpdateProduct($id: ID!, $input: UpdateProductInput!) {
  updateProduct(id: $id, input: $input) {
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
    }
    isActive
    createdAt
    updatedAt
  }
}

# Delete product
mutation DeleteProduct($id: ID!) {
  deleteProduct(id: $id)
}
```

#### Order Mutations
```graphql
# Create order
mutation CreateOrder($input: CreateOrderInput!) {
  createOrder(input: $input) {
    id
    customerId
    customer {
      id
      firstName
      lastName
      email
    }
    orderNumber
    status
    totalAmount
    shippingAddress
    billingAddress
    notes
    orderDate
    orderItems {
      id
      productId
      product {
        id
        name
        price
      }
      quantity
      unitPrice
      totalPrice
    }
    createdAt
    updatedAt
  }
}

# Update order
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

# Delete order
mutation DeleteOrder($id: ID!) {
  deleteOrder(id: $id)
}
```

### Input Types

#### CreateUserInput
```graphql
input CreateUserInput {
  name: String!
  email: String!
}
```

#### UpdateUserInput
```graphql
input UpdateUserInput {
  name: String
  email: String
}
```

#### CreateCustomerInput
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
```

#### UpdateCustomerInput
```graphql
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

#### CreateCategoryInput
```graphql
input CreateCategoryInput {
  name: String!
  description: String
  parentId: ID
}
```

#### UpdateCategoryInput
```graphql
input UpdateCategoryInput {
  name: String
  description: String
  parentId: ID
}
```

#### CreateProductInput
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
```

#### UpdateProductInput
```graphql
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

#### CreateOrderInput
```graphql
input CreateOrderInput {
  customerId: ID!
  shippingAddress: String!
  billingAddress: String!
  notes: String
  orderItems: [CreateOrderItemInput!]!
}
```

#### CreateOrderItemInput
```graphql
input CreateOrderItemInput {
  productId: ID!
  quantity: Int!
}
```

#### UpdateOrderInput
```graphql
input UpdateOrderInput {
  status: OrderStatus
  shippingAddress: String
  billingAddress: String
  notes: String
  shippedDate: Time
  deliveredDate: Time
}
```

## Error Handling

### HTTP Status Codes

- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters or body
- `401 Unauthorized`: Authentication required or invalid token
- `403 Forbidden`: Access denied
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

### Error Response Format

**REST API:**
```json
{
  "error": "Error message describing what went wrong"
}
```

**GraphQL API:**
```json
{
  "errors": [
    {
      "message": "Error message",
      "locations": [
        {
          "line": 2,
          "column": 3
        }
      ],
      "path": ["fieldName"]
    }
  ],
  "data": null
}
```

## Examples

### REST API Examples

#### Create a Product
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{
    "name": "MacBook Pro",
    "description": "Apple MacBook Pro 16-inch",
    "sku": "MBP16-512GB",
    "price": 2499.99,
    "stock": 25,
    "category_id": "electronics-category-uuid",
    "is_active": true
  }'
```

#### Get Products with Filtering
```bash
curl "http://localhost:8080/api/products?limit=20&offset=0&category_id=electronics-uuid&is_active=true"
```

#### Create an Order
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{
    "customer_id": "customer-uuid",
    "shipping_address": "123 Main St, New York, NY 10001",
    "billing_address": "123 Main St, New York, NY 10001",
    "notes": "Please deliver during business hours",
    "order_items": [
      {
        "product_id": "product-uuid",
        "quantity": 2
      }
    ]
  }'
```

### GraphQL Examples

#### Query Products with Category Information
```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query GetProducts($limit: Int, $categoryId: ID) { products(limit: $limit, categoryId: $categoryId) { id name description price stock category { id name } isActive } }",
    "variables": {
      "limit": 10,
      "categoryId": "electronics-category-uuid"
    }
  }'
```

#### Create Order with Items
```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation CreateOrder($input: CreateOrderInput!) { createOrder(input: $input) { id orderNumber status totalAmount customer { firstName lastName email } orderItems { product { name } quantity unitPrice totalPrice } } }",
    "variables": {
      "input": {
        "customerId": "customer-uuid",
        "shippingAddress": "123 Main St, New York, NY 10001",
        "billingAddress": "123 Main St, New York, NY 10001",
        "notes": "Handle with care",
        "orderItems": [
          {
            "productId": "product-uuid",
            "quantity": 1
          }
        ]
      }
    }
  }'
```

### OIDC Authentication Flow Example

#### Step 1: Get Authorization URL
```bash
curl http://localhost:8080/auth/oidc/login
```

#### Step 2: Handle Callback (after user authorizes)
```bash
curl "http://localhost:8080/auth/oidc/callback?code=auth-code&state=state-value"
```

#### Step 3: Use Token for Customer Operations
```bash
curl -H "Authorization: Bearer <oidc-jwt-token>" \
  http://localhost:8080/api/customer/profile
```

## Additional Resources

- **GraphQL Playground**: Visit `http://localhost:8080/graphql/playground` for interactive GraphQL exploration
- **OIDC Setup**: See [OIDC_SETUP.md](./OIDC_SETUP.md) for detailed OIDC configuration
- **Project README**: See [README.md](./README.md) for general project information

## Support

For questions or issues with the API, please refer to the project documentation or contact the development team.