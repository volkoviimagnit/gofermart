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
	users map[string]model.User
}

func NewUserRepositoryMem() IUserRepository {
	return &UserRepositoryMem{users: map[string]model.User{}}
}

func (u *UserRepositoryMem) Insert(user model.User) error {
	user.SetId(randStr(10))
	u.users[user.Login()] = user
	return nil
}

func (u *UserRepositoryMem) GetOneByCredentials(login string, password string) (*model.User, error) {
	logrus.Debugf("UserRepositoryMem.users %+v", u.users)

	user, isExist := u.users[login]
	if !isExist {
		return nil, errors.New("несуществующая пара логин/пароль")
	}
	if user.Password() != password {
		return nil, errors.New("неверная пара логин/пароль")
	}
	return &user, nil
}

func (u *UserRepositoryMem) Update(user model.User) error {
	u.users[user.Login()] = user
	return nil
}

func randStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	// Read b number of numbers
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
