package security

import (
	"errors"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/helpers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
)

type Authenticator struct {
	userRepository repository.IUserRepository
}

func (a Authenticator) Authenticate(request *http.Request) (IPassport, error) {
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

func (a Authenticator) CreateAuthenticatedToken() string {
	return helpers.RandomString(10)
}

func (a Authenticator) GetUserRepository() repository.IUserRepository {
	//TODO implement me
	panic("implement me")
}

func NewAuthenticator(userRepository repository.IUserRepository) *Authenticator {
	return &Authenticator{
		userRepository: userRepository,
	}
}
