package models

import "github.com/Hanasou/news_feed/go/common"

// User represents a user in the system
type User struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (user *User) ToJson() (string, error) {
	return common.ToJson(user)
}

func (user *User) GetId() (string, error) {
	return user.Id, nil
}
