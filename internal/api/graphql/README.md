# GraphQL API Documentation

This directory contains the complete GraphQL API implementation for the SIL Backend Assessment project, built with [gqlgen](https://gqlgen.com/) following GraphQL best practices.

## 📁 Directory Structure

```
internal/api/graphql/
├── README.md                    # This documentation
├── router.go                    # GraphQL router with middleware and configuration
├── graph/
│   ├── schema.graphqls         # GraphQL schema definition
│   ├── generated.go            # Generated GraphQL server code
│   ├── resolver.go             # Generated resolver interfaces
│   └── model/
│       └── models_gen.go       # Generated GraphQL models
└── resolvers/
    ├── resolver.go             # Main resolver with dependency injection
    └── schema.resolvers.go     # Resolver implementations
```

## 🚀 Features

### Core Functionality
- **Complete CRUD Operations** for all entities (Users, Customers, Categories, Products, Orders)
- **Relationship Queries** with automatic resolution
- **Advanced Filtering** and pagination
- **Search Capabilities** across entities
- **Order Management** with status tracking
- **Statistics and Analytics** endpoints

### Technical Features
- **Type Safety** with Go type system integration
- **Automatic Code Generation** via gqlgen
- **Dependency Injection** for clean architecture
- **Comprehensive Error Handling** with proper GraphQL errors
- **Request Logging** and performance monitoring
- **CORS Support** for web clients
- **Query Caching** for improved performance
- **Automatic Persisted Queries** support

## 📊 Schema Overview

### Core Types
- **User**: System users with JWT authentication
- **Customer**: Customers with OIDC authentication  
- **Category**: Hierarchical product categories
- **Product**: Products with inventory tracking
- **Order**: Orders with status management
- **OrderItem**: Individual items within orders

### Relationships
```graphql
Customer -> Orders -> OrderItems -> Products -> Categories
                                      ↓
                                  Categories (hierarchical)
```

## 🔧 Configuration

### Router Configuration
The GraphQL router is configured with:
- **Multiple Transports**: WebSocket, GET, POST, Multipart
- **Query Cache**: LRU cache with 1000 entries
- **Extensions**: Introspection, Automatic Persisted Queries
- **Error Handling**: Custom error presenter with logging
- **Panic Recovery**: Graceful error handling

### Performance Optimizations
- **Query Caching**: Reduces parsing overhead
- **Automatic Persisted Queries**: Reduces bandwidth
- **Efficient Pagination**: Limit/offset with sensible defaults
- **Relationship Loading**: Optimized N+1 query prevention

## 📝 API Usage

### Endpoints
- **GraphQL Endpoint**: `POST /graphql`
- **GraphQL Playground**: `GET /graphql/playground`
- **Schema Introspection**: `GET /graphql/schema`
- **Health Check**: `GET /graphql/health`

### Authentication
The GraphQL API supports the same authentication methods as the REST API:
- **JWT Authentication**: For user operations
- **OIDC Authentication**: For customer operations

### Example Queries

#### Basic Query
```graphql
query GetUsers {
  users(pagination: { limit: 10, offset: 0 }) {
    id
    name
    email
    createdAt
  }
}
```

#### Relationship Query
```graphql
query GetProductsWithCategories {
  products(pagination: { limit: 5 }) {
    id
    name
    price
    stock
    category {
      id
      name
      parent {
        name
      }
    }
  }
}
```

#### Complex Query with Filtering
```graphql
query GetOrdersByCustomer($customerId: ID!) {
  ordersByCustomer(
    customerId: $customerId
    pagination: { limit: 20 }
  ) {
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
```

#### Mutation Example
```graphql
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
```

### Variables Example
```json
{
  "input": {
    "name": "iPhone 15 Pro",
    "description": "Latest iPhone model",
    "sku": "IPHONE15PRO-256GB",
    "price": 1199.99,
    "stock": 50,
    "categoryId": "category-uuid-here",
    "isActive": true
  }
}
```

## 🔍 Advanced Features

### Pagination
All list queries support pagination:
```graphql
{
  products(pagination: { limit: 20, offset: 40 }) {
    id
    name
  }
}
```

### Filtering
Products and orders support advanced filtering:
```graphql
{
  products(
    filter: {
      categoryId: "category-id"
      isActive: true
      minPrice: 10.0
      maxPrice: 100.0
    }
  ) {
    id
    name
    price
  }
}
```

### Search
Search across multiple fields:
```graphql
{
  searchProducts(query: "iPhone", pagination: { limit: 10 }) {
    id
    name
    description
    sku
  }
}
```

### Order Management
Specialized order operations:
```graphql
mutation ShipOrder($orderId: ID!) {
  shipOrder(id: $orderId) {
    id
    status
    shippedDate
  }
}
```

### Statistics
Get business insights:
```graphql
{
  orderStats {
    totalOrders
    totalRevenue
    averageOrderValue
    ordersByStatus {
      status
      count
    }
  }
}
```

## 🛠️ Development

### Code Generation
The GraphQL code is automatically generated using gqlgen:

```bash
# Generate GraphQL code
go run github.com/99designs/gqlgen generate

# Or use the Makefile
make generate-graphql
```

### Adding New Fields
1. Update `schema.graphqls` with new fields
2. Run code generation
3. Implement the resolver methods
4. Update tests and documentation

### Adding New Types
1. Define the type in `schema.graphqls`
2. Add to `gqlgen.yml` model mapping if needed
3. Run code generation
4. Implement all required resolvers

## 🧪 Testing

### GraphQL Playground
Access the interactive playground at `http://localhost:8080/graphql/playground` to:
- Explore the schema
- Test queries and mutations
- View documentation
- Debug issues

### Example Test Queries

#### Health Check
```graphql
query {
  users(pagination: { limit: 1 }) {
    id
  }
}
```

#### Create and Query Flow
```graphql
# 1. Create a category
mutation {
  createCategory(input: { name: "Electronics", description: "Electronic devices" }) {
    id
    name
  }
}

# 2. Create a product
mutation {
  createProduct(input: {
    name: "Laptop"
    sku: "LAPTOP-001"
    price: 999.99
    stock: 10
    categoryId: "category-id-from-step-1"
  }) {
    id
    name
    category {
      name
    }
  }
}

# 3. Query the relationship
query {
  categories {
    name
    products {
      name
      price
    }
  }
}
```

## 🔒 Security

### Input Validation
- All inputs are validated at the GraphQL layer
- UUID validation for ID fields
- Business logic validation in services

### Error Handling
- Sensitive information is not exposed in errors
- Proper error logging for debugging
- GraphQL-compliant error responses

### Rate Limiting
Consider implementing rate limiting for production:
- Query complexity analysis
- Depth limiting
- Request rate limiting

## 📈 Performance

### Query Optimization
- Use pagination for large datasets
- Avoid deep nesting when possible
- Use specific field selection

### Caching
- Query result caching
- Automatic persisted queries
- CDN caching for static schema

### Monitoring
- Request logging with timing
- Error rate monitoring
- Query complexity tracking

## 🚀 Deployment

### Environment Variables
```bash
# GraphQL specific settings
GRAPHQL_PLAYGROUND_ENABLED=true
GRAPHQL_INTROSPECTION_ENABLED=true
GRAPHQL_QUERY_CACHE_SIZE=1000
```

### Production Considerations
- Disable playground in production
- Enable query complexity limiting
- Set up proper monitoring
- Configure CORS appropriately

## 📚 Resources

### Documentation
- [GraphQL Specification](https://spec.graphql.org/)
- [gqlgen Documentation](https://gqlgen.com/)
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)

### Tools
- [GraphQL Playground](https://github.com/graphql/graphql-playground)
- [GraphiQL](https://github.com/graphql/graphiql)
- [Apollo Studio](https://studio.apollographql.com/)

## 🤝 Contributing

### Code Style
- Follow Go conventions
- Use meaningful resolver names
- Add comprehensive logging
- Handle errors gracefully

### Schema Design
- Use descriptive field names
- Add comprehensive documentation
- Follow GraphQL naming conventions
- Design for client needs

### Testing
- Test all resolvers
- Validate error handling
- Test relationship loading
- Performance test with large datasets