package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type IUserBalanceRepository interface {
	Insert(row model.UserBalance) error
	FinOneByUserId(userId string) (*model.UserBalance, error)
	Update(row model.UserBalance) error
	Upset(row model.UserBalance) error
}
