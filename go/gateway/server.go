package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Hanasou/news_feed/go/common/auth"
	"github.com/Hanasou/news_feed/go/gateway/graph"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

// JWTMiddleware validates JWT tokens for GraphQL requests
// TODO: Send request to user service to validate JWT
// and fetch user information.
func JWTMiddleware(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for introspection queries and playground
			if r.URL.Path == "/" || isIntrospectionQuery(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// For GraphQL, we might want to allow unauthenticated requests
				// and handle authorization at the resolver level
				next.ServeHTTP(w, r)
				return
			}

			token, err := auth.ExtractTokenFromHeader(authHeader)
			if err != nil {
				http.Error(w, "Invalid authorization header: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Validate token
			claims, err := jwtService.ValidateAccessToken(token)
			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Add user info to request context
			ctx := auth.WithUserContext(r.Context(), claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// isIntrospectionQuery checks if the request is a GraphQL introspection query
func isIntrospectionQuery(r *http.Request) bool {
	if r.Method != "POST" {
		return false
	}

	// Simple check for introspection - you might want to make this more sophisticated
	contentType := r.Header.Get("Content-Type")
	return strings.Contains(contentType, "application/json") &&
		r.Header.Get("X-GraphQL-Introspection") == "true"
}

func main() {
	// Initialize JWT service
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		// In production, always use environment variables for secrets
		secretKey = "your-super-secret-key-min-32-chars-long"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	jwtService := auth.NewJWTService(secretKey, "news-feed-gateway")

	startGraphQlServer(jwtService)
	log.Println("GraphQL server started on port", defaultPort)
	select {} // Keep the server running
}

func startGraphQlServer(jwtService *auth.JWTService) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// TODO: Initialize clients here.
	// Get client types from config file
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Create JWT middleware
	jwtMiddleware := JWTMiddleware(jwtService)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", jwtMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
