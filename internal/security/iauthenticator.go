package security

import (
	"net/http"
)

type IAuthenticator interface {
	Authenticate(request *http.Request) (IPassport, error)
	SetAuthenticatedToken(rw http.ResponseWriter, token string)
	CreateAuthenticatedToken() string
}
