package grpc

// grpc UserClient is a client for interacting with the User service over gRPC.

import (
	"context"
	"log"

	"github.com/Hanasou/news_feed/go/common/auth"
	"github.com/Hanasou/news_feed/go/common/common_models"
	"github.com/Hanasou/news_feed/go/common/common_models/responses"
	"github.com/Hanasou/news_feed/go/common/grpc/userpb"
)

type GrpcUserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(client userpb.UserServiceClient) *GrpcUserClient {
	return &GrpcUserClient{client: client}
}

func (c *GrpcUserClient) CreateUser(ctx context.Context, user *common_models.User) (*responses.CreateUserResponse, error) {

	newUser := &userpb.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role.String(),
	}

	req := &userpb.CreateUserRequest{
		User: newUser,
	}
	grpcUserResponse, err := c.client.CreateUser(ctx, req)
	if err != nil {
		log.Println("Error in CreateUser from User service: ", err)
		return nil, err
	}

	return &responses.CreateUserResponse{
		Response: grpcUserResponse.Response,
		User:     user,
	}, nil
}

func (c *GrpcUserClient) AuthenticateUser(ctx context.Context, identifier, password string) (*responses.AuthUserResponse, error) {
	req := &userpb.AuthenticateUserRequest{Identifier: identifier, Password: password}
	grpcAuthResponse, err := c.client.AuthenticateUser(ctx, req)
	if err != nil {
		log.Println("Error in AuthenticateUser from User service: ", err)
		return nil, err
	}
	user := &common_models.User{
		ID:       grpcAuthResponse.User.ID,
		Username: grpcAuthResponse.User.Username,
		Email:    grpcAuthResponse.User.Email,
		Password: grpcAuthResponse.User.Password,
		Role:     common_models.RoleFromString(grpcAuthResponse.User.Role),
	}
	return &responses.AuthUserResponse{
		TokenPair: &auth.TokenPair{
			AccessToken:  grpcAuthResponse.AccessToken,
			RefreshToken: grpcAuthResponse.RefreshToken,
			ExpiresIn:    grpcAuthResponse.ExpiresTimestamp,
			TokenType:    grpcAuthResponse.TokenType,
		},
		User: user,
	}, nil
}
