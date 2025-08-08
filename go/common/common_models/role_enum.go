package common_models

import "fmt"

type Role string

const (
	Admin   Role = "admin"
	Default Role = "default"
)

func (r Role) IsValid() bool {
	switch r {
	case Admin, Default:
		return true
	default:
		return false
	}
}

func (r Role) String() string {
	return string(r)
}

func RoleFromString(roleStr string) (Role, error) {
	switch roleStr {
	case "admin":
		return Admin, nil
	case "default":
		return Default, nil
	default:
		return "", fmt.Errorf("invalid role: %s", roleStr)
	}
}
