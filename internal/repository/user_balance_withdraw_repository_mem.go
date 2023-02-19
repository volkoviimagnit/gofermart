package repository

import (
	"sync"

	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserBalanceWithdrawRepositoryMem struct {
	userWithdraws map[string][]model.UserBalanceWithdraw
	mutex         *sync.RWMutex
}

func NewUserBalanceWithdrawRepositoryMem() IUserBalanceWithdrawRepository {
	return &UserBalanceWithdrawRepositoryMem{
		userWithdraws: make(map[string][]model.UserBalanceWithdraw, 0),
		mutex:         &sync.RWMutex{},
	}
}

func (r *UserBalanceWithdrawRepositoryMem) Insert(row model.UserBalanceWithdraw) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	userID := row.UserID
	if _, isExist := r.userWithdraws[userID]; !isExist {
		r.userWithdraws[userID] = make([]model.UserBalanceWithdraw, 0)
	}
	r.userWithdraws[userID] = append(r.userWithdraws[userID], row)
	return nil
}

func (r *UserBalanceWithdrawRepositoryMem) FindByUserID(userID string) ([]model.UserBalanceWithdraw, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if _, isExist := r.userWithdraws[userID]; !isExist {
		return make([]model.UserBalanceWithdraw, 0), nil
	}
	return r.userWithdraws[userID], nil
}

func (r *UserBalanceWithdrawRepositoryMem) SumWithdrawByUserID(userID string) (float64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	sum := 0.0
	for _, userWithDraw := range r.userWithdraws[userID] {
		sum += userWithDraw.Sum
	}
	return sum, nil
}
