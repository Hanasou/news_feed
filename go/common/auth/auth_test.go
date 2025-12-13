package auth

import (
	"testing"
	"time"

	"github.com/Hanasou/news_feed/go/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTAuthentication(t *testing.T) {
	// Create JWT service
	secretKey := "your-super-secret-key-min-32-chars-long"
	jwtService := NewJWTService(secretKey, "news-feed-service")

	// Create test user
	user := &models.User{
		ID:       "user123",
		Username: "john_doe",
		Email:    "john@example.com",
		Role:     models.Default,
	}

	// Hash password
	hashedPassword, err := HashPassword("mypassword123")
	require.NoError(t, err)
	user.Password = hashedPassword

	t.Run("Generate and validate token pair", func(t *testing.T) {
		// Generate token pair
		tokens, err := jwtService.GenerateTokenPair(user)
		require.NoError(t, err)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.NotEmpty(t, tokens.RefreshToken)
		assert.Equal(t, "Bearer", tokens.TokenType)
		assert.Greater(t, tokens.ExpiresIn, int64(0))

		// Validate access token
		claims, err := jwtService.ValidateAccessToken(tokens.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Username, claims.Username)
		assert.Equal(t, user.Email, claims.Email)
		assert.Equal(t, user.Role, claims.Role)

		// Validate refresh token
		userID, err := jwtService.ValidateRefreshToken(tokens.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, user.ID, userID)
	})

	t.Run("Password validation", func(t *testing.T) {
		// Valid password
		err := ValidatePassword(user.Password, "mypassword123")
		assert.NoError(t, err)

		// Invalid password
		err = ValidatePassword(user.Password, "wrongpassword")
		assert.Error(t, err)
	})

	t.Run("Token header extraction", func(t *testing.T) {
		// Valid header
		token, err := ExtractTokenFromHeader("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
		require.NoError(t, err)
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", token)

		// Invalid header format
		_, err = ExtractTokenFromHeader("InvalidHeader")
		assert.Error(t, err)

		// Empty header
		_, err = ExtractTokenFromHeader("")
		assert.Error(t, err)
	})

	t.Run("Expired token validation", func(t *testing.T) {
		// Create JWT service with very short expiry
		shortExpiryService := &JWTService{
			secretKey:     []byte(secretKey),
			issuer:        "news-feed-service",
			accessExpiry:  time.Millisecond * 10, // 10ms
			refreshExpiry: time.Hour * 24 * 7,
		}

		// Generate token
		tokens, err := shortExpiryService.GenerateTokenPair(user)
		require.NoError(t, err)

		// Wait for token to expire
		time.Sleep(time.Millisecond * 20)

		// Try to validate expired token
		_, err = shortExpiryService.ValidateAccessToken(tokens.AccessToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token has invalid claims")
	})

	t.Run("Invalid token validation", func(t *testing.T) {
		// Invalid token format
		_, err := jwtService.ValidateAccessToken("invalid.token.format")
		assert.Error(t, err)

		// Empty token
		_, err = jwtService.ValidateAccessToken("")
		assert.Error(t, err)

		// Token with wrong signature
		_, err = jwtService.ValidateAccessToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
		assert.Error(t, err)
	})
}

func TestGenerateSecureRandomString(t *testing.T) {
	// Test random string generation
	random1, err := GenerateSecureRandomString(32)
	require.NoError(t, err)
	assert.Len(t, random1, 64) // hex encoding doubles the length

	random2, err := GenerateSecureRandomString(32)
	require.NoError(t, err)
	assert.Len(t, random2, 64)

	// Should generate different strings
	assert.NotEqual(t, random1, random2)
}

func BenchmarkPasswordHashing(b *testing.B) {
	password := "testpassword123"

	b.Run("HashPassword", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := HashPassword(password)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Pre-hash for validation benchmark
	hashedPassword, _ := HashPassword(password)

	b.Run("ValidatePassword", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := ValidatePassword(hashedPassword, password)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkJWTOperations(b *testing.B) {
	secretKey := "your-super-secret-key-min-32-chars-long"
	jwtService := NewJWTService(secretKey, "news-feed-service")

	user := &models.User{
		ID:       "user123",
		Username: "john_doe",
		Email:    "john@example.com",
		Role:     models.Default,
	}

	b.Run("GenerateTokenPair", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := jwtService.GenerateTokenPair(user)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Pre-generate token for validation benchmark
	tokens, _ := jwtService.GenerateTokenPair(user)

	b.Run("ValidateAccessToken", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := jwtService.ValidateAccessToken(tokens.AccessToken)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
