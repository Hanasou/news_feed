package grpc

// implement grpc server for todo service

import (
	"context"
	"log"

	"github.com/Hanasou/news_feed/go/common/grpc/todopb"
	"github.com/Hanasou/news_feed/go/todo/core"
	"github.com/Hanasou/news_feed/go/todo/models"
)

type TodoServer struct {
	todopb.UnimplementedTodoServiceServer
	service *core.TodoService
}

func NewTodoServer(service *core.TodoService) *TodoServer {
	return &TodoServer{service: service}
}

func (s *TodoServer) CreateTodo(ctx context.Context, req *todopb.CreateTodoRequest) (*todopb.CreateTodoResponse, error) {
	todo := &models.Todo{
		Id:     req.Todo.GetId(),
		Text:   req.Todo.GetText(),
		Done:   req.Todo.GetDone(),
		UserId: req.Todo.GetUserId(),
	}

	err := s.service.CreateTodo(todo)
	if err != nil {
		log.Printf("Failed to create todo: %v", err)
		return nil, err
	}

	return &todopb.CreateTodoResponse{Response: "Todo created successfully"}, nil
}

func (s *TodoServer) GetTodos(ctx context.Context, req *todopb.GetTodosRequest) (*todopb.GetTodosResponse, error) {
	todos, err := s.service.GetTodos(req.GetUserId())
	if err != nil {
		log.Printf("Failed to get todos: %v", err)
		return nil, err
	}

	var todoList []*todopb.Todo
	for _, todo := range todos {
		todoList = append(todoList, &todopb.Todo{
			Id:     todo.Id,
			Text:   todo.Text,
			Done:   todo.Done,
			UserId: todo.UserId,
		})
	}

	return &todopb.GetTodosResponse{Todos: todoList}, nil
}
