# JWT Authentication for News Feed Service

This package provides comprehensive JWT (JSON Web Token) authentication for the News Feed service, including secure password hashing, token generation/validation, and middleware for HTTP and gRPC services.

## Features

- **Secure Password Hashing**: Uses bcrypt for password hashing with salt
- **JWT Token Pairs**: Access tokens (short-lived) and refresh tokens (long-lived)  
- **Context Management**: User information stored in request context
- **HTTP Middleware**: Ready-to-use middleware for HTTP services
- **gRPC Interceptors**: Authentication interceptors for gRPC services
- **Comprehensive Testing**: Unit tests and benchmarks included

## Quick Start

### 1. Basic Setup

```go
import "github.com/Hanasou/news_feed/go/user/auth"

// Create JWT service
secretKey := "your-super-secret-key-min-32-chars-long"
jwtService := auth.NewJWTService(secretKey, "news-feed-service")
```

### 2. User Registration

```go
// Hash password
hashedPassword, err := auth.HashPassword("userpassword")
if err != nil {
    log.Fatal(err)
}

// Create user
user := &auth.User{
    ID:       "user123",
    Username: "john_doe", 
    Email:    "john@example.com",
    Password: hashedPassword,
    Role:     "user",
}
```

### 3. User Login

```go
// Validate password
err := auth.ValidatePassword(user.Password, "userpassword")
if err != nil {
    // Invalid password
    return
}

// Generate token pair
tokens, err := jwtService.GenerateTokenPair(user)
if err != nil {
    log.Fatal(err)
}

// tokens.AccessToken - use for API calls
// tokens.RefreshToken - use to get new access tokens
// tokens.ExpiresIn - access token expiry time in seconds
```

### 4. Token Validation

```go
// Validate access token
claims, err := jwtService.ValidateAccessToken(accessToken)
if err != nil {
    // Invalid or expired token
    return
}

// Extract user info
userID := claims.UserID
username := claims.Username
email := claims.Email
role := claims.Role
```

## HTTP Middleware Usage

```go
func main() {
    jwtService := auth.NewJWTService(secretKey, "news-feed-service")
    
    // Create middleware
    jwtMiddleware := auth.JWTMiddleware(jwtService)
    
    // Protected route
    http.Handle("/api/protected", jwtMiddleware(protectedHandler))
    
    // Public route
    http.Handle("/api/login", loginHandler)
    
    http.ListenAndServe(":8080", nil)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    // Extract user info from context
    userID, _ := auth.GetUserIDFromContext(r.Context())
    username, _ := auth.GetUsernameFromContext(r.Context())
    
    fmt.Fprintf(w, "Hello %s (ID: %s)", username, userID)
}
```

## gRPC Interceptor Usage

```go
func main() {
    jwtService := auth.NewJWTService(secretKey, "news-feed-service")
    jwtInterceptor := auth.NewJWTInterceptor(jwtService)
    
    // Create gRPC server with JWT middleware
    server := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            jwtInterceptor.UnaryServerInterceptor(),
        ),
        grpc.ChainStreamInterceptor(
            jwtInterceptor.StreamServerInterceptor(),
        ),
    )
    
    // Register services
    pb.RegisterUserServiceServer(server, &userServiceImpl{})
    pb.RegisterTodoServiceServer(server, &todoServiceImpl{})
    
    // Start server
    lis, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    server.Serve(lis)
}
```

## Client Usage

### HTTP Client

```bash
# Login to get tokens
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "john_doe", "password": "userpassword"}'

# Use access token for protected endpoints
curl -X GET http://localhost:8080/api/protected \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### gRPC Client

```go
// Add token to gRPC metadata
ctx := context.Background()
ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+accessToken)

// Make authenticated gRPC call
response, err := client.GetTodos(ctx, &pb.GetTodosRequest{})
```

## Configuration

The JWT service can be configured with custom expiry times:

```go
jwtService := &auth.JWTService{
    // Custom configuration
}
```

Default settings:
- **Access Token Expiry**: 1 hour
- **Refresh Token Expiry**: 7 days
- **Hash Algorithm**: HMAC-SHA256
- **Password Hashing**: bcrypt with cost 12

## Security Best Practices

1. **Secret Key**: Use a strong secret key (minimum 32 characters)
2. **HTTPS Only**: Always use HTTPS in production
3. **Token Storage**: Store tokens securely on client side
4. **Token Rotation**: Implement token refresh mechanism
5. **Logout**: Implement proper logout with token invalidation

## Testing

```bash
# Run tests
go test ./user/auth -v

# Run benchmarks
go test ./user/auth -bench=. -benchmem
```

## Dependencies

- `github.com/golang-jwt/jwt/v5` - JWT token handling
- `golang.org/x/crypto/bcrypt` - Password hashing
- `google.golang.org/grpc` - gRPC interceptors
- `github.com/stretchr/testify` - Testing (dev dependency)

## API Reference

### Core Types

- `JWTService` - Main service for JWT operations
- `User` - User information structure
- `Claims` - JWT token claims
- `TokenPair` - Access and refresh token pair

### Key Functions

- `NewJWTService(secretKey, issuer string) *JWTService`
- `GenerateTokenPair(user *User) (*TokenPair, error)`
- `ValidateAccessToken(token string) (*Claims, error)`
- `ValidateRefreshToken(token string) (string, error)`
- `HashPassword(password string) (string, error)`
- `ValidatePassword(hashedPassword, password string) error`

### Context Utilities

- `WithUserContext(ctx context.Context, claims *Claims) context.Context`
- `GetUserIDFromContext(ctx context.Context) (string, error)`
- `GetUsernameFromContext(ctx context.Context) (string, error)`
- `GetUserEmailFromContext(ctx context.Context) (string, error)`
- `GetUserRoleFromContext(ctx context.Context) (string, error)`

## Performance

Based on benchmarks:
- **Token Generation**: ~6.7μs per operation
- **Token Validation**: ~6.0μs per operation  
- **Password Hashing**: ~45ms per operation (intentionally slow for security)
- **Password Validation**: ~44ms per operation

The JWT operations are very fast, while password operations are intentionally slower for security against brute force attacks.
