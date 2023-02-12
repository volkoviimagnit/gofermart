package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type IUserBalanceWithdrawRepository interface {
	Insert(row model.UserBalanceWithdraw) error
	FindByUserID(userID string) ([]model.UserBalanceWithdraw, error)
	SumWithdrawByUserID(userID string) (float64, error)
}
