package security

type IUser interface {
	GetLogin() string
	GetID() string
}
