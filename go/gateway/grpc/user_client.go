package grpc

// grpc UserClient is a client for interacting with the User service over gRPC.

import (
	"context"

	"github.com/Hanasou/news_feed/go/common/grpc/userpb"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(client userpb.UserServiceClient) *UserClient {
	return &UserClient{client: client}
}

func (c *UserClient) CreateUser(ctx context.Context, user *userpb.User) (*userpb.CreateUserResponse, error) {
	req := &userpb.CreateUserRequest{User: user}
	return c.client.CreateUser(ctx, req)
}

func (c *UserClient) AuthenticateUser(ctx context.Context, identifier, password string) (*userpb.AuthenticateUserResponse, error) {
	req := &userpb.AuthenticateUserRequest{Identifier: identifier, Password: password}
	return c.client.AuthenticateUser(ctx, req)
}
