package responses

import (
	"github.com/Hanasou/news_feed/go/common/auth"
	"github.com/Hanasou/news_feed/go/common/common_models"
)

type AuthUserResponse struct {
	TokenPair *auth.TokenPair
	User      *common_models.User
}

type CreateUserResponse struct {
	Response string
	User     *common_models.User
}
