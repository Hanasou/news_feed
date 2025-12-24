package server

import (
	"context"

	"github.com/Hanasou/news_feed/go/common/models"
	"github.com/Hanasou/news_feed/go/common/models/responses"
)

// This interface may not actually be necessary
type UserServer interface {
	CreateUserRequest(context.Context, *models.User) (*responses.CreateUserResponse, error)
	AuthenticateUserRequest(context.Context, string, string) (*responses.AuthUserResponse, error)
}
