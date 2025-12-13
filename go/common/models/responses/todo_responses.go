package responses

import "github.com/Hanasou/news_feed/go/common/models"

type CreateTodoResponse struct {
	Response string
	TodoID   string
}

type GetTodosResponse struct {
	Response string
	Todos    []models.Todo
}
