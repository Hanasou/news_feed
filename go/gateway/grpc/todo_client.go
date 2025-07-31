package grpc

// grpc TodoClient is a client for interacting with the Todo service over gRPC.

import (
	"context"

	"github.com/Hanasou/news_feed/go/common/grpc/todopb"
)

type TodoClient struct {
	client todopb.TodoServiceClient
}

func NewTodoClient(client todopb.TodoServiceClient) *TodoClient {
	return &TodoClient{client: client}
}

func (c *TodoClient) CreateTodo(ctx context.Context, todo *todopb.Todo) (*todopb.CreateTodoResponse, error) {
	req := &todopb.CreateTodoRequest{Todo: todo}
	return c.client.CreateTodo(ctx, req)
}

func (c *TodoClient) GetTodos(ctx context.Context, userId string) (*todopb.GetTodosResponse, error) {
	req := &todopb.GetTodosRequest{UserId: userId}
	return c.client.GetTodos(ctx, req)
}
