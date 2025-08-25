package server

import (
	"context"

	"github.com/Hanasou/news_feed/go/common/common_models"
	"github.com/Hanasou/news_feed/go/common/common_models/responses"
)

type UserServer interface {
	CreateUserRequest(context.Context, *common_models.User) (*responses.CreateUserResponse, error)
	AuthenticateUserRequest(context.Context, string, string) (*responses.AuthUserResponse, error)
}
