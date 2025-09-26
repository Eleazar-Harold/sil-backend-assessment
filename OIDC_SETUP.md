# OpenID Connect Authentication Setup

This document explains how to set up and use OpenID Connect authentication for customers in the SIL Backend Assessment application.

## Overview

The application now supports both traditional JWT-based authentication for users and OpenID Connect authentication for customers. This allows customers to authenticate using external identity providers like Google, Microsoft, Auth0, etc.

## Configuration

### 1. Update config.yaml

```yaml
oidc:
  enabled: true
  provider_url: https://accounts.google.com  # Your OIDC provider URL
  client_id: your-google-client-id          # Your OIDC client ID
  client_secret: your-google-client-secret  # Your OIDC client secret
  redirect_url: http://localhost:8080/auth/oidc/callback
  scopes:
    - openid
    - profile
    - email
```

### 2. Environment Variables (Alternative)

You can also configure OIDC using environment variables:

```bash
export OIDC_ENABLED=true
export OIDC_PROVIDER_URL=https://accounts.google.com
export OIDC_CLIENT_ID=your-google-client-id
export OIDC_CLIENT_SECRET=your-google-client-secret
export OIDC_REDIRECT_URL=http://localhost:8080/auth/oidc/callback
```

## Supported OIDC Providers

### Google OAuth 2.0
- Provider URL: `https://accounts.google.com`
- Required scopes: `openid`, `profile`, `email`

### Microsoft Azure AD
- Provider URL: `https://login.microsoftonline.com/{tenant-id}/v2.0`
- Required scopes: `openid`, `profile`, `email`

### Auth0
- Provider URL: `https://your-domain.auth0.com`
- Required scopes: `openid`, `profile`, `email`

### Generic OIDC Provider
- Provider URL: Your OIDC provider's well-known configuration URL
- Required scopes: `openid`, `profile`, `email`

## API Endpoints

### Authentication Flow

#### 1. Get Authorization URL
```http
GET /auth/oidc/login
```

Response:
```json
{
  "auth_url": "https://accounts.google.com/oauth/authorize?...",
  "state": "random-state-string"
}
```

#### 2. Handle Callback
```http
GET /auth/oidc/callback?code=...&state=...
```

Response:
```json
{
  "access_token": "jwt-access-token",
  "refresh_token": "jwt-refresh-token",
  "customer": {
    "id": "customer-uuid",
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

#### 3. Validate Token
```http
GET /auth/oidc/validate
Authorization: Bearer <jwt-token>
```

Response:
```json
{
  "sub": "google-user-id",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "given_name": "John",
  "family_name": "Doe",
  "picture": "https://..."
}
```

### Using the token in GraphQL Playground

1. Open `http://localhost:8080/graphql/playground`.
2. Click the HTTP Headers pane and add:

```
{
  "Authorization": "Bearer <your-oidc-or-jwt-token>"
}
```

3. Run queries/mutations, for example:

```
query { products(pagination: { limit: 5 }) { id name } }
```

### Customer Profile Management

#### Get Profile
```http
GET /api/customer/profile
Authorization: Bearer <jwt-token>
```

#### Update Profile
```http
PUT /api/customer/profile
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "address": "123 Main St",
  "city": "New York",
  "state": "NY",
  "zip_code": "10001",
  "country": "USA"
}
```

#### Delete Account
```http
DELETE /api/customer/account
Authorization: Bearer <jwt-token>
```

## Frontend Integration

### 1. Redirect to OIDC Provider

```javascript
// Get authorization URL
const response = await fetch('/auth/oidc/login');
const { auth_url, state } = await response.json();

// Store state for validation
localStorage.setItem('oidc_state', state);

// Redirect to OIDC provider
window.location.href = auth_url;
```

### 2. Handle Callback

```javascript
// In your callback page
const urlParams = new URLSearchParams(window.location.search);
const code = urlParams.get('code');
const state = urlParams.get('state');

// Validate state
const storedState = localStorage.getItem('oidc_state');
if (state !== storedState) {
  console.error('Invalid state parameter');
  return;
}

// Exchange code for tokens
const response = await fetch(`/auth/oidc/callback?code=${code}&state=${state}`);
const authData = await response.json();

// Store tokens
localStorage.setItem('access_token', authData.access_token);
localStorage.setItem('refresh_token', authData.refresh_token);

// Redirect to main application
window.location.href = '/dashboard';
```

### 3. Make Authenticated Requests

```javascript
// Include token in requests
const token = localStorage.getItem('access_token');

const response = await fetch('/api/customer/profile', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});

const profile = await response.json();
```

## Security Considerations

1. **State Parameter**: Always validate the state parameter to prevent CSRF attacks
2. **Token Storage**: Store tokens securely (consider using httpOnly cookies for production)
3. **HTTPS**: Always use HTTPS in production
4. **Token Expiry**: Implement token refresh logic
5. **Scope Validation**: Only request necessary scopes from the OIDC provider

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK`: Successful authentication/profile operations
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Invalid or expired tokens
- `500 Internal Server Error`: Server-side errors

Error responses include descriptive messages:

```json
{
  "error": "Invalid token: token expired"
}
```

## Testing

### 1. Test OIDC Flow

```bash
# Start the server
go run cmd/server/main_with_oidc.go

# Test authorization URL generation
curl http://localhost:8080/auth/oidc/login

# Test token validation (replace with actual token)
curl -H "Authorization: Bearer <token>" http://localhost:8080/auth/oidc/validate
```

### 2. Test Customer Endpoints

```bash
# Get customer profile
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/customer/profile

# Update customer profile
curl -X PUT -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"first_name": "Jane"}' \
  http://localhost:8080/api/customer/profile
```

## Migration from Traditional Auth

If you're migrating from traditional username/password authentication:

1. Keep existing user authentication for admin/management functions
2. Use OIDC authentication for customer-facing features
3. Both authentication methods can coexist in the same application
4. The middleware automatically handles both JWT and OIDC tokens

## Troubleshooting

### Common Issues

1. **"OIDC provider not configured"**: Check that `oidc.enabled` is `true` and all required fields are set
2. **"Invalid client credentials"**: Verify your OIDC client ID and secret
3. **"Invalid redirect URI"**: Ensure the redirect URL matches what's configured in your OIDC provider
4. **"Token validation failed"**: Check that the token is properly formatted and not expired

### Debug Mode

Enable debug logging by setting the log level to `debug` in your configuration:

```yaml
logging:
  level: debug
```

This will provide detailed information about the OIDC flow and token validation process.
