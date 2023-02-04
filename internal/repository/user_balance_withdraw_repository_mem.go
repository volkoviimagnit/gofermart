package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type UserBalanceWithdrawRepositoryMem struct {
	userWithdraws map[string][]model.UserBalanceWithdraw
}

func (r UserBalanceWithdrawRepositoryMem) Insert(row model.UserBalanceWithdraw) error {
	userId := row.GetUserId()
	if _, isExist := r.userWithdraws[userId]; !isExist {
		r.userWithdraws[userId] = make([]model.UserBalanceWithdraw, 0)
	}
	r.userWithdraws[userId] = append(r.userWithdraws[userId], row)
	return nil
}

func (r UserBalanceWithdrawRepositoryMem) FindByUserId(userId string) ([]model.UserBalanceWithdraw, error) {
	if _, isExist := r.userWithdraws[userId]; !isExist {
		return make([]model.UserBalanceWithdraw, 0), nil
	}
	return r.userWithdraws[userId], nil
}

func NewUserBalanceWithdrawRepositoryMem() IUserBalanceWithdrawRepository {
	return &UserBalanceWithdrawRepositoryMem{
		userWithdraws: make(map[string][]model.UserBalanceWithdraw, 0),
	}
}
