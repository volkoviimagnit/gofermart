package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type IUserOrderRepository interface {
	Insert(row model.UserOrder) error
	Update(row model.UserOrder) error
	FindByUserID(userID string) ([]model.UserOrder, error)
	FindOneByNumber(number string) (*model.UserOrder, error)
	SumAccrualByUserID(userID string) (float64, error)
}
