package security

type Passport struct {
	user IUser
}

func (p *Passport) GetUser() IUser {
	return p.user
}

func NewPassport(user IUser) *Passport {
	return &Passport{user: user}
}
