package repository

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserRepositoryMem struct {
	loginUsers map[string]model.User
	// todo глупое хранилище юзеров по токенам
	tokenUsers map[string]model.User
}

func NewUserRepositoryMem() IUserRepository {
	return &UserRepositoryMem{
		loginUsers: map[string]model.User{},
		tokenUsers: map[string]model.User{},
	}
}

func (u *UserRepositoryMem) Insert(user model.User) error {
	user.SetId(randStr(10))
	u.addLoginUser(user)
	u.addTokenUser(user)
	return nil
}

func (u *UserRepositoryMem) FindOneByCredentials(login string, password string) (*model.User, error) {
	logrus.Debugf("UserRepositoryMem.loginUsers %+v", u.loginUsers)

	user, isExist := u.loginUsers[login]
	if !isExist {
		return nil, errors.New("несуществующая пара логин/пароль")
	}
	if user.Password() != password {
		return nil, errors.New("неверная пара логин/пароль")
	}
	return &user, nil
}

func (u *UserRepositoryMem) FindOneByLogin(login string) (*model.User, error) {
	logrus.Debugf("UserRepositoryMem.loginUsers %+v", u.loginUsers)

	user, isExist := u.loginUsers[login]
	if !isExist {
		return nil, nil
	}
	return &user, nil
}

func (u *UserRepositoryMem) FindOneByToken(token string) (*model.User, error) {
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
	u.loginUsers[user.Login()] = user
}
func (u *UserRepositoryMem) addTokenUser(user model.User) {
	userToken := user.Token()
	if len(userToken) > 0 {
		u.tokenUsers[userToken] = user
	}
}

func randStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	// Read b number of numbers
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
