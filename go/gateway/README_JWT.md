# JWT Authentication Middleware for GraphQL Gateway

This implementation adds JWT authentication to your GraphQL server with the following features:

## Features

- ✅ JWT token validation on all GraphQL requests
- ✅ Automatic context injection with user information
- ✅ Configurable public endpoints (playground, introspection)
- ✅ Flexible authentication requirements per resolver
- ✅ Role-based access control support

## How It Works

### 1. Middleware Layer

The `JWTMiddleware` function wraps your GraphQL handler and:
- Extracts JWT tokens from the `Authorization` header
- Validates tokens using the JWT service
- Injects user information into the request context
- Allows unauthenticated requests to pass through (resolver-level auth)

### 2. Resolver Integration

Resolvers can access authenticated user information using:

```go
// Get current user claims
claims, err := auth.GetClaimsFromContext(ctx)
if err != nil {
    return nil, fmt.Errorf("authentication required: %w", err)
}

// Get specific user information
userID, err := auth.GetUserIDFromContext(ctx)
username, err := auth.GetUsernameFromContext(ctx)
userRole, err := auth.GetUserRoleFromContext(ctx)
```

### 3. Authorization Patterns

#### Require Authentication
```go
func (r *queryResolver) ProtectedQuery(ctx context.Context) (*Model, error) {
    claims, err := auth.GetClaimsFromContext(ctx)
    if err != nil {
        return nil, fmt.Errorf("authentication required: %w", err)
    }
    // ... rest of resolver logic
}
```

#### Require Specific Role
```go
func (r *queryResolver) AdminOnlyQuery(ctx context.Context) (*Model, error) {
    userRole, err := auth.GetUserRoleFromContext(ctx)
    if err != nil || userRole != "admin" {
        return nil, fmt.Errorf("admin access required")
    }
    // ... rest of resolver logic
}
```

#### Resource Ownership Check
```go
func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, input NewTodo) (*Todo, error) {
    claims, err := auth.GetClaimsFromContext(ctx)
    if err != nil {
        return nil, fmt.Errorf("authentication required: %w", err)
    }
    
    // Check if user owns the todo or is admin
    // ... fetch todo from database
    if todo.UserID != claims.UserID {
        userRole, _ := auth.GetUserRoleFromContext(ctx)
        if userRole != "admin" {
            return nil, fmt.Errorf("access denied")
        }
    }
    // ... rest of resolver logic
}
```

## Configuration

### Environment Variables

- `JWT_SECRET`: Secret key for JWT signing (required in production)
- `PORT`: Server port (default: 8080)

### JWT Service Settings

The JWT service is configured with:
- Access token expiry: 15 minutes
- Refresh token expiry: 7 days
- Issuer: "news-feed-gateway"

## Usage Examples

### 1. Start the Server

```bash
cd /home/rzhang/projects/news_feed/go/gateway
JWT_SECRET="your-super-secret-key-min-32-chars-long" go run .
```

### 2. Authenticate

```graphql
mutation {
  authenticateUser(input: { 
    identifier: "username", 
    password: "password" 
  }) {
    accessToken
    refreshToken
    user {
      id
      name
      email
      role
    }
  }
}
```

### 3. Make Authenticated Requests

Include the access token in the Authorization header:

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"query": "{ todos { id text done } }"}' \
  http://localhost:8080/query
```

### 4. GraphQL Playground

The playground at `http://localhost:8080/` includes an HTTP Headers section where you can add:

```json
{
  "Authorization": "Bearer YOUR_ACCESS_TOKEN"
}
```

## Security Considerations

1. **Secret Key**: Always use a strong, randomly generated secret key in production
2. **HTTPS**: Use HTTPS in production to protect tokens in transit
3. **Token Storage**: Store tokens securely on the client side
4. **Token Expiry**: Short-lived access tokens with refresh token rotation
5. **Rate Limiting**: Consider adding rate limiting to prevent brute force attacks

## Testing

Run the test script to verify JWT authentication:

```bash
./test_jwt.sh
```

This script tests:
- Unauthenticated access
- User authentication
- Authenticated queries
- Protected mutations
- Role-based access control

## Error Handling

The middleware returns appropriate HTTP status codes:
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Valid token but insufficient permissions (handled at resolver level)

## Integration with Existing Services

To integrate with your existing gRPC services:

1. Pass JWT service to resolvers via dependency injection
2. Extract user context and forward to gRPC calls
3. Use the same JWT validation in gRPC middleware (see `grpc_middleware.go`)

## Next Steps

1. Implement actual user lookup and password validation
2. Add refresh token rotation
3. Integrate with your existing user and todo services
4. Add logging and monitoring
5. Implement rate limiting
6. Add CORS configuration for browser clients
