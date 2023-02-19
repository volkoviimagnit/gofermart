package model

type User struct {
	ID       string
	Login    string
	Password string
	Token    string
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetLogin() string {
	return u.Login
}
