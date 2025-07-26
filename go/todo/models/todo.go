package models

import "github.com/Hanasou/news_feed/go/common"

type Todo struct {
	Id     string `json:"id,omitempty"`
	Text   string `json:"text,omitempty"`
	Done   bool   `json:"done,omitempty"`
	UserId string `json:"user_id,omitempty"`
}

func (todo *Todo) ToJson() (string, error) {
	return common.ToJson(todo)
}

func (todo *Todo) GetId() (string, error) {
	return todo.Id, nil
}
