package model

type User struct {
	id       string
	login    string
	password string
	token    string
}

func (u *User) Id() string {
	return u.id
}

func (u *User) SetId(id string) {
	u.id = id
}

func (u *User) Login() string {
	return u.login
}

func (u *User) SetLogin(login string) {
	u.login = login
}

func (u *User) Password() string {
	return u.password
}

func (u *User) SetPassword(password string) {
	u.password = password
}

func (u *User) Token() string {
	return u.token
}

func (u *User) SetToken(token string) {
	u.token = token
}
