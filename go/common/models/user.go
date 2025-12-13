package models

import (
	"fmt"

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

func (user *User) ToMap() (map[string]any, error) {
	return common.ToMap(user)
}

func (user *User) GetID() (string, error) {
	return user.ID, nil
}

func (user *User) GetField(field string) (any, error) {
	switch field {
	case "id":
		return user.ID, nil
	case "username":
		return user.Username, nil
	case "email":
		return user.Email, nil
	case "password":
		return user.Password, nil
	case "role":
		return user.Role, nil
	default:
		return nil, fmt.Errorf("field %s not found", field)
	}
}
