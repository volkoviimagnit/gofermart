package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type IUserOrderRepository interface {
	Insert(row model.UserOrder) error
	FindByUserId(userId string) ([]model.UserOrder, error)
	FindOneByNumber(number string) (*model.UserOrder, error)
	IsExist(userId string, number string) (bool, error)
}
