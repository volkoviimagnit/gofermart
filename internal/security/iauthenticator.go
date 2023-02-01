package security

import (
	"net/http"
)

type IAuthenticator interface {
	Authenticate(request *http.Request) (IPassport, error)
	CreateAuthenticatedToken() string
}
