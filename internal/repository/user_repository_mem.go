package repository

import (
	"errors"
	"sync"

	"github.com/volkoviimagnit/gofermart/internal/helpers"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserRepositoryMem struct {
	loginUsers map[string]model.User
	// todo глупое хранилище юзеров по токенам
	tokenUsers map[string]model.User
	mutex      *sync.RWMutex
}

func NewUserRepositoryMem() IUserRepository {
	return &UserRepositoryMem{
		loginUsers: map[string]model.User{},
		tokenUsers: map[string]model.User{},
		mutex:      &sync.RWMutex{},
	}
}

func (u *UserRepositoryMem) Insert(user model.User) error {
	user.ID = helpers.RandomString(10)
	u.addLoginUser(user)
	u.addTokenUser(user)
	return nil
}

// FindOneByCredentials TODO:  добавить шифрование пароля и сверку с шифром, а не самим паролем
func (u *UserRepositoryMem) FindOneByCredentials(login string, password string) (*model.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	user, isExist := u.loginUsers[login]
	if !isExist {
		return nil, errors.New("несуществующая пара логин/пароль")
	}
	if user.Password != password {
		return nil, errors.New("неверная пара логин/пароль")
	}
	return &user, nil
}

func (u *UserRepositoryMem) FindOneByLogin(login string) (*model.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	user, isExist := u.loginUsers[login]
	if !isExist {
		return nil, nil
	}
	return &user, nil
}

func (u *UserRepositoryMem) FindOneByToken(token string) (*model.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	user, isExist := u.tokenUsers[token]
	if !isExist {
		return nil, nil
	}
	return &user, nil
}

func (u *UserRepositoryMem) Update(user model.User) error {
	u.addLoginUser(user)
	u.addTokenUser(user)
	return nil
}

func (u *UserRepositoryMem) addLoginUser(user model.User) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.loginUsers[user.GetLogin()] = user
}
func (u *UserRepositoryMem) addTokenUser(user model.User) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	userToken := user.Token
	if len(userToken) > 0 {
		u.tokenUsers[userToken] = user
	}
}
