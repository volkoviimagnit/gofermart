package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type IUserRepository interface {
	Insert(user model.User) error
	GetOneByCredentials(login string, password string) (*model.User, error)
	GetOneByLogin(login string) (*model.User, error)
	GetOneByToken(token string) (*model.User, error)
	Update(user model.User) error
}
