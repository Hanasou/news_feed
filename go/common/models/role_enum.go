package models

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

func RoleFromString(roleStr string) Role {
	switch roleStr {
	case "admin":
		return Admin
	case "default":
		return Default
	default:
		return Default
	}
}
