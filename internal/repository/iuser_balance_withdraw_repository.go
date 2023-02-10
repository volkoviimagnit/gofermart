package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type IUserBalanceWithdrawRepository interface {
	Insert(row model.UserBalanceWithdraw) error
	FindByUserId(userId string) ([]model.UserBalanceWithdraw, error)
	SumWithdrawByUserId(userId string) float64
}
