package main

import (
	"context"
	"errors"

	"github.com/Hanasou/news_feed/go/user/auth"
)

// GraphQL context helpers for authentication

// GetCurrentUser extracts the current user's claims from GraphQL context
func GetCurrentUser(ctx context.Context) (*auth.Claims, error) {
	return auth.GetClaimsFromContext(ctx)
}

// GetCurrentUserID extracts the current user's ID from GraphQL context
func GetCurrentUserID(ctx context.Context) (string, error) {
	return auth.GetUserIDFromContext(ctx)
}

// RequireAuth ensures a user is authenticated for a GraphQL operation
func RequireAuth(ctx context.Context) (*auth.Claims, error) {
	claims, err := GetCurrentUser(ctx)
	if err != nil {
		return nil, errors.New("authentication required")
	}
	return claims, nil
}

// RequireRole ensures a user has a specific role for a GraphQL operation
func RequireRole(ctx context.Context, requiredRole string) (*auth.Claims, error) {
	claims, err := RequireAuth(ctx)
	if err != nil {
		return nil, err
	}

	userRole, err := auth.GetUserRoleFromContext(ctx)
	if err != nil {
		return nil, errors.New("user role not found")
	}

	if userRole != requiredRole {
		return nil, errors.New("insufficient permissions")
	}

	return claims, nil
}

// IsOwnerOrAdmin checks if the current user is the owner of a resource or an admin
func IsOwnerOrAdmin(ctx context.Context, resourceOwnerID string) (bool, error) {
	claims, err := GetCurrentUser(ctx)
	if err != nil {
		return false, err
	}

	// Check if user is the owner
	if claims.UserID == resourceOwnerID {
		return true, nil
	}

	// Check if user is an admin
	userRole, err := auth.GetUserRoleFromContext(ctx)
	if err != nil {
		return false, err
	}

	return userRole == "admin", nil
}
