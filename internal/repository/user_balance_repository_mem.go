package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type UserBalanceRepositoryMem struct {
}

func NewUserBalanceRepositoryMem() IUserBalanceRepository {
	return &UserBalanceRepositoryMem{}
}

func (u UserBalanceRepositoryMem) Insert(row model.UserBalance) error {
	//TODO implement me
	panic("implement me")
}

func (u UserBalanceRepositoryMem) FinOneByUserId(userId string) (model.UserBalance, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserBalanceRepositoryMem) Update(row model.UserBalance) error {
	//TODO implement me
	panic("implement me")
}
