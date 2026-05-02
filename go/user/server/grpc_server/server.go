package grpc_server

import (
	"context"
	"log"

	"github.com/Hanasou/news_feed/go/common/grpc/userpb"
	"github.com/Hanasou/news_feed/go/common/models"
	"github.com/Hanasou/news_feed/go/user/core"
)

type GrpcUserServer struct {
	userpb.UnimplementedUserServiceServer
	service *core.UserService
}

func NewGrpcUserServer(service *core.UserService) *GrpcUserServer {
	return &GrpcUserServer{
		service: service,
	}
}

func (s *GrpcUserServer) CreateUser(ctx context.Context, request *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	responseMessage := "User failed to get created"
	user := &models.User{
		ID:       request.User.Id,
		Username: request.User.Username,
		Email:    request.User.Email,
		Password: request.User.Password,
		Role:     models.RoleFromString(request.User.Role),
	}
	err := s.service.CreateUser(user)
	if err != nil {
		log.Println(responseMessage)
		return &userpb.CreateUserResponse{
			Response: responseMessage,
		}, err
	}

	responseMessage = "User succesfully added"
	return &userpb.CreateUserResponse{
		Response: responseMessage,
	}, nil
}

func (s *GrpcUserServer) AuthenticateUser(ctx context.Context, request *userpb.AuthenticateUserRequest) (*userpb.AuthenticateUserResponse, error) {
	responseMessage := "User failed to get authenticated"
	tokenPair, user, err := s.service.AuthenticateUser(request.Identifier, request.Password)
	if err != nil {
		log.Println(responseMessage)
		return nil, err
	}
	response := &userpb.AuthenticateUserResponse{
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		ExpiresTimestamp: tokenPair.ExpiresIn,
		TokenType:        tokenPair.TokenType,
		User: &userpb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Role:     user.Role.String(),
		},
	}
	return response, nil
}

func (s *GrpcUserServer) GetUsers(ctx context.Context, request *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	users, err := s.service.GetUsers(request.IdFilter, request.NameFilter, request.EmailFilter, request.RoleFilter)
	if err != nil {
		log.Println("Failed to get users")
		return nil, err
	}
	response := &userpb.GetUsersResponse{
		Response: "Users retrieved successfully",
		Users:    make([]*userpb.User, len(users)),
	}
	for i, user := range users {
		response.Users[i] = &userpb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Role:     user.Role.String(),
		}
	}
	return response, nil
}

func (s *GrpcUserServer) mustEmbedUnimplementedUserServiceServer() {}
