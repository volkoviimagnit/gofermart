package security

import (
	"errors"
	"math/rand"
	"net/http"

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
	return a.randomString(10)
}

func (a Authenticator) randomString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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
