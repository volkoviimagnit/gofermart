package security

import (
	"net/http"
)

type IAuthenticator interface {
	Authenticate(request *http.Request) (IPassport, error)
	RenderAuthenticatedToken(rw http.ResponseWriter, token string)
	SetAuthenticatedToken(r *http.Request, token string)
	CreateAuthenticatedToken() string
}
