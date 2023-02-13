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

func (u *User) SetID(id string) {
	u.ID = id
}

func (u *User) GetLogin() string {
	return u.Login
}

func (u *User) SetLogin(login string) {
	u.Login = login
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

func (u *User) GetToken() string {
	return u.Token
}

func (u *User) SetToken(token string) {
	u.Token = token
}
