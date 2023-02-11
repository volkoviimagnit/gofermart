package security

type IUser interface {
	GetLogin() string
	GetId() string
}
