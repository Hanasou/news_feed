package clients

import (
	"context"

	"github.com/Hanasou/news_feed/go/common/models"
	"github.com/Hanasou/news_feed/go/common/models/responses"
)

type UserClient interface {
	CreateUser(context.Context, *models.User) (*responses.CreateUserResponse, error)
	AuthenticateUser(context.Context, string, string) (*responses.AuthUserResponse, error)
}

type TodoClient interface {
}
