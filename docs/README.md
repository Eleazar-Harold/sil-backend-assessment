# API Documentation

This directory contains the interactive API documentation for the SIL Backend Assessment project.

## Available Documentation

### REST API Documentation

#### OpenAPI/Swagger Specification
- **File**: `swagger.json`
- **Endpoint**: `http://localhost:8080/swagger.json`
- **Description**: Complete OpenAPI 3.0.3 specification for all REST endpoints

#### Swagger UI
- **File**: `swagger-ui.html`
- **Endpoint**: `http://localhost:8080/swagger-ui.html`
- **Description**: Interactive Swagger UI for testing REST API endpoints
- **Features**:
  - Try-it-out functionality
  - Request/response examples
  - Schema validation
  - Authentication support

#### ReDoc
- **File**: `redoc.html`
- **Endpoint**: `http://localhost:8080/redoc.html`
- **Description**: Beautiful, responsive API documentation using ReDoc
- **Features**:
  - Clean, professional layout
  - Searchable documentation
  - Code samples
  - Three-panel design

### GraphQL API Documentation

#### GraphQL Documentation
- **File**: `graphql-docs.html`
- **Endpoint**: `http://localhost:8080/graphql-docs.html`
- **Description**: Comprehensive GraphQL API documentation
- **Features**:
  - Complete schema overview
  - Query and mutation examples
  - Type definitions
  - Interactive examples

#### GraphQL Playground
- **Endpoint**: `http://localhost:8080/graphql/playground`
- **Description**: Interactive GraphQL IDE (served by gqlgen)
- **Features**:
  - Query execution
  - Schema exploration
  - Auto-completion
  - Query history

## Documentation Endpoints

### Main Documentation Index
- **Endpoint**: `http://localhost:8080/docs`
- **Description**: JSON response with links to all documentation resources

### Convenience Routes
- `http://localhost:8080/docs/rest` → Swagger UI
- `http://localhost:8080/docs/graphql` → GraphQL documentation
- `http://localhost:8080/docs/redoc` → ReDoc documentation

## File Structure

```
docs/
├── README.md              # This file
├── swagger.json           # OpenAPI 3.0.3 specification
├── swagger-ui.html        # Swagger UI interface
├── redoc.html            # ReDoc interface
└── graphql-docs.html     # GraphQL documentation
```

## Features

### REST API Documentation Features
- **Complete Coverage**: All REST endpoints documented
- **Authentication**: JWT and OIDC authentication examples
- **Request/Response Examples**: Real-world examples for all operations
- **Error Handling**: Comprehensive error response documentation
- **Pagination**: Query parameter documentation for paginated endpoints
- **Filtering**: Documentation for filtering and search parameters

### GraphQL API Documentation Features
- **Schema Introspection**: Complete GraphQL schema documentation
- **Query Examples**: Real-world query examples with variables
- **Mutation Examples**: Complete mutation examples with input types
- **Relationship Queries**: Examples of nested queries with relationships
- **Type Definitions**: Complete type system documentation

## Usage

### For Developers
1. **API Exploration**: Use Swagger UI or GraphQL Playground for interactive exploration
2. **Integration**: Use the OpenAPI spec for code generation
3. **Testing**: Use the try-it-out features for endpoint testing

### For Documentation
1. **Reference**: Use ReDoc for clean, printable documentation
2. **Examples**: Copy examples from the documentation for implementation
3. **Schema**: Reference the complete schema definitions

## Updating Documentation

The documentation is automatically served from these static files. To update:

1. **REST API**: Update `swagger.json` with new endpoints or changes
2. **GraphQL**: Update `graphql-docs.html` with new schema information
3. **UI Updates**: Modify the HTML files to update styling or functionality

## Integration with Project

The documentation is integrated into the main application through:

- **Handler**: `internal/api/docs/handler.go`
- **Routes**: Registered in `cmd/server/main.go`
- **Static Files**: Served directly from this directory

## External Dependencies

### Swagger UI
- **Version**: 5.9.0
- **CDN**: unpkg.com
- **License**: Apache 2.0

### ReDoc
- **Version**: 2.1.3
- **CDN**: jsdelivr.net
- **License**: MIT

## Browser Compatibility

All documentation interfaces are compatible with:
- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+

## Security Considerations

- All documentation is served with CORS headers for development
- No sensitive information is exposed in the documentation
- Authentication examples use placeholder tokens
- Production deployments should consider access restrictions