package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type IUserRepository interface {
	Insert(user model.User) error
	FindOneByCredentials(login string, password string) (*model.User, error)
	FindOneByLogin(login string) (*model.User, error)
	FindOneByToken(token string) (*model.User, error)
	Update(user model.User) error
}
