package models

import (
	"fmt"

	"github.com/Hanasou/news_feed/go/common"
)

type Todo struct {
	Id     string `json:"id,omitempty"`
	Text   string `json:"text,omitempty"`
	Done   bool   `json:"done,omitempty"`
	UserId string `json:"user_id,omitempty"`
}

func (todo *Todo) ToJson() (string, error) {
	return common.ToJson(todo)
}

func (todo *Todo) GetID() (string, error) {
	return todo.Id, nil
}

func (todo *Todo) ToMap() (map[string]any, error) {
	return common.ToMap(todo)
}

func (todo *Todo) GetField(field string) (any, error) {
	switch field {
	case "id":
		return todo.Id, nil
	case "text":
		return todo.Text, nil
	case "done":
		return todo.Done, nil
	case "user_id":
		return todo.UserId, nil
	default:
		return nil, fmt.Errorf("field %s not found", field)
	}
}
