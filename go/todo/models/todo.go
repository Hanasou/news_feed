package models

import "github.com/Hanasou/news_feed/go/common/commodels"

type Todo struct {
	Id      string          `json:"id,omitempty"`
	Text    string          `json:"text,omitempty"`
	Done    bool            `json:"done,omitempty"`
	User_id *commodels.User `json:"user_id,omitempty"`
}

func (todo *Todo) ToJson() (string, error) {
	return commodels.ToJson(todo)
}

func (todo *Todo) GetId() (string, error) {
	return todo.Id, nil
}
