package security

type IUser interface {
	Login() string
	Id() string
}
