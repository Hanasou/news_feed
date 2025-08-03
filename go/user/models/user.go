package models

import (
	"github.com/Hanasou/news_feed/go/common"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // hashed
	Role     Role   `json:"role"`     // User role (e.g., "admin
}

func (user *User) ToJson() (string, error) {
	return common.ToJson(user)
}

func (user *User) GetId() (string, error) {
	return user.ID, nil
}
