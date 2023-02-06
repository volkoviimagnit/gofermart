package security

import (
	"errors"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/helpers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
)

type AuthenticatorHeader struct {
	userRepository repository.IUserRepository
}

func (a *AuthenticatorHeader) Authenticate(request *http.Request) (IPassport, error) {
	accessToken := request.Header.Get("Authorization")

	user, errToking := a.userRepository.FindOneByToken(accessToken)
	if errToking != nil {
		return nil, errors.New("не удалось проверить токен")
	}
	if user == nil {
		return nil, nil
	}
	passport := NewPassport(user)
	return passport, nil

}

func (a *AuthenticatorHeader) CreateAuthenticatedToken() string {
	return helpers.RandomOrderNumber()
}

func (a *AuthenticatorHeader) GetUserRepository() repository.IUserRepository {
	//TODO implement me
	panic("implement me")
}

func (a *AuthenticatorHeader) RenderAuthenticatedToken(rw http.ResponseWriter, token string) {
	rw.Header().Set("Authorization", token)
}

func (a *AuthenticatorHeader) SetAuthenticatedToken(r *http.Request, token string) {
	r.Header.Set("Authorization", token)
}

func NewAuthenticator(userRepository repository.IUserRepository) *AuthenticatorHeader {
	return &AuthenticatorHeader{
		userRepository: userRepository,
	}
}
