package clients

import (
	"context"

	"github.com/Hanasou/news_feed/go/common/common_models"
	"github.com/Hanasou/news_feed/go/common/common_models/responses"
)

type UserClient interface {
	CreateUser(context.Context, *common_models.User) (*responses.CreateUserResponse, error)
	AuthenticateUser(context.Context, string, string) (*responses.AuthUserResponse, error)
}

type TodoClient interface {
}
