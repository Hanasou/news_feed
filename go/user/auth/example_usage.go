package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Example usage of the JWT authentication system
func ExampleUsage() {
	// 1. Create JWT service
	secretKey := "your-super-secret-key-min-32-chars-long"
	jwtService := NewJWTService(secretKey, "news-feed-service")

	// 2. Create a user (normally from database)
	user := &User{
		ID:       "user123",
		Username: "john_doe",
		Email:    "john@example.com",
		Role:     "user",
	}

	// Hash password when creating user
	hashedPassword, err := HashPassword("mypassword123")
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}
	user.Password = hashedPassword

	// 3. Generate token pair (login)
	tokens, err := jwtService.GenerateTokenPair(user)
	if err != nil {
		log.Fatal("Failed to generate tokens:", err)
	}

	fmt.Printf("Access Token: %s\n", tokens.AccessToken)
	fmt.Printf("Refresh Token: %s\n", tokens.RefreshToken)
	fmt.Printf("Expires In: %d seconds\n", tokens.ExpiresIn)

	// 4. Validate access token (middleware usage)
	claims, err := jwtService.ValidateAccessToken(tokens.AccessToken)
	if err != nil {
		log.Fatal("Failed to validate token:", err)
	}

	fmt.Printf("Token valid for user: %s (%s)\n", claims.Username, claims.Email)

	// 5. Add user info to context
	ctx := context.Background()
	ctx = WithUserContext(ctx, claims)

	// 6. Extract user info from context (in handlers)
	userID, _ := GetUserIDFromContext(ctx)
	username, _ := GetUsernameFromContext(ctx)
	fmt.Printf("Context User: %s (ID: %s)\n", username, userID)

	// 7. Validate password (login)
	err = ValidatePassword(user.Password, "mypassword123")
	if err != nil {
		fmt.Println("Invalid password")
	} else {
		fmt.Println("Password valid!")
	}

	// 8. Refresh token usage
	userIDFromRefresh, err := jwtService.ValidateRefreshToken(tokens.RefreshToken)
	if err != nil {
		log.Fatal("Failed to validate refresh token:", err)
	}
	fmt.Printf("Refresh token valid for user ID: %s\n", userIDFromRefresh)
}

// HTTP Middleware example
func JWTMiddleware(jwtService *JWTService) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			token, err := ExtractTokenFromHeader(authHeader)
			if err != nil {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Validate token
			claims, err := jwtService.ValidateAccessToken(token)
			if err != nil {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Add user info to request context
			ctx := WithUserContext(r.Context(), claims)
			r = r.WithContext(ctx)

			// Call next handler
			next.ServeHTTP(w, r)
		}
	}
}

// Login handler example
func LoginHandler(jwtService *JWTService, userService UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// Parse request body
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Find user by username
		user, err := userService.GetByUsername(loginRequest.Username)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Validate password
		if err := ValidatePassword(user.Password, loginRequest.Password); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate token pair
		tokens, err := jwtService.GenerateTokenPair(user)
		if err != nil {
			http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
			return
		}

		// Return tokens
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens)
	}
}

// UserService interface (implement this)
type UserService interface {
	GetByUsername(username string) (*User, error)
	GetByID(id string) (*User, error)
	Create(user *User) error
}
