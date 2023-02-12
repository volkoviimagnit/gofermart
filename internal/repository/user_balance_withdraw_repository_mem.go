package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type UserBalanceWithdrawRepositoryMem struct {
	userWithdraws map[string][]model.UserBalanceWithdraw
}

func (r UserBalanceWithdrawRepositoryMem) Insert(row model.UserBalanceWithdraw) error {
	userID := row.GetUserId()
	if _, isExist := r.userWithdraws[userID]; !isExist {
		r.userWithdraws[userID] = make([]model.UserBalanceWithdraw, 0)
	}
	r.userWithdraws[userID] = append(r.userWithdraws[userID], row)
	return nil
}

func (r UserBalanceWithdrawRepositoryMem) FindByUserID(userId string) ([]model.UserBalanceWithdraw, error) {
	if _, isExist := r.userWithdraws[userId]; !isExist {
		return make([]model.UserBalanceWithdraw, 0), nil
	}
	return r.userWithdraws[userId], nil
}

func (r UserBalanceWithdrawRepositoryMem) SumWithdrawByUserID(userID string) (float64, error) {
	sum := 0.0
	for _, userWithDraw := range r.userWithdraws[userID] {
		sum += userWithDraw.GetSum()
	}
	return sum, nil
}

func NewUserBalanceWithdrawRepositoryMem() IUserBalanceWithdrawRepository {
	return &UserBalanceWithdrawRepositoryMem{
		userWithdraws: make(map[string][]model.UserBalanceWithdraw, 0),
	}
}
