package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type UserBalanceRepositoryMem struct {
	userBalance map[string]model.UserBalance
}

func NewUserBalanceRepositoryMem() IUserBalanceRepository {
	return &UserBalanceRepositoryMem{
		userBalance: make(map[string]model.UserBalance, 0),
	}
}

func (u *UserBalanceRepositoryMem) Insert(row model.UserBalance) error {
	u.userBalance[row.UserID] = row
	return nil
}

func (u *UserBalanceRepositoryMem) FinOneByUserID(userID string) (*model.UserBalance, error) {
	if _, isExist := u.userBalance[userID]; !isExist {
		return nil, nil
	}
	row := u.userBalance[userID]
	return &row, nil
}

func (u *UserBalanceRepositoryMem) Update(row model.UserBalance) error {
	u.userBalance[row.UserID] = row
	return nil
}

func (u *UserBalanceRepositoryMem) Upset(row model.UserBalance) error {
	u.userBalance[row.UserID] = row
	return nil
}
