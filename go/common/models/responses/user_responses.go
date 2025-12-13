package responses

import (
	"github.com/Hanasou/news_feed/go/common/auth"
	"github.com/Hanasou/news_feed/go/common/models"
)

type AuthUserResponse struct {
	TokenPair *auth.TokenPair
	User      *models.User
}

type CreateUserResponse struct {
	Response string
	User     *models.User
}
