package commodels

type User struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (user *User) ToJson() (string, error) {
	return ToJson(user)
}

func (user *User) GetId() (string, error) {
	return user.Id, nil
}
