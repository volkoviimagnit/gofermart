package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

// IUserBalanceRepository TODO: избавить от методов  Insert and Update
type IUserBalanceRepository interface {
	Insert(row model.UserBalance) error
	FinOneByUserID(userID string) (*model.UserBalance, error)
	Update(row model.UserBalance) error
	Upset(row model.UserBalance) error
}
