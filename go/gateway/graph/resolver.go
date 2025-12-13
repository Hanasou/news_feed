package graph

import (
	"github.com/Hanasou/news_feed/go/gateway/clients"
	"github.com/Hanasou/news_feed/go/gateway/config"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserClient clients.UserClient
	TodoClient clients.TodoClient
	Config     *config.GatewayConfig
}
