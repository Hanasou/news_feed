package grpc_clients

// grpc TodoClient is a client for interacting with the Todo service over gRPC.

import (
	"context"

	"github.com/Hanasou/news_feed/go/common/grpc/todopb"
)

type GrpcTodoClient struct {
	client todopb.TodoServiceClient
}

func NewTodoClient(client todopb.TodoServiceClient) *GrpcTodoClient {
	return &GrpcTodoClient{client: client}
}

func (c *GrpcTodoClient) CreateTodo(ctx context.Context, todo *todopb.Todo) (*todopb.CreateTodoResponse, error) {
	req := &todopb.CreateTodoRequest{Todo: todo}
	return c.client.CreateTodo(ctx, req)
}

func (c *GrpcTodoClient) GetTodos(ctx context.Context, userId string) (*todopb.GetTodosResponse, error) {
	req := &todopb.GetTodosRequest{UserId: userId}
	return c.client.GetTodos(ctx, req)
}
