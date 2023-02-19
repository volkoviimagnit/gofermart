package repository

import (
	"sync"

	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserBalanceRepositoryMem struct {
	userBalance map[string]model.UserBalance
	mutex       *sync.RWMutex
}

func NewUserBalanceRepositoryMem() IUserBalanceRepository {
	return &UserBalanceRepositoryMem{
		userBalance: make(map[string]model.UserBalance, 0),
		mutex:       &sync.RWMutex{},
	}
}

func (u *UserBalanceRepositoryMem) Insert(row model.UserBalance) error {
	return u.Upset(row)
}

func (u *UserBalanceRepositoryMem) FinOneByUserID(userID string) (*model.UserBalance, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	if _, isExist := u.userBalance[userID]; !isExist {
		return nil, nil
	}
	row := u.userBalance[userID]
	return &row, nil
}

func (u *UserBalanceRepositoryMem) Update(row model.UserBalance) error {
	return u.Upset(row)
}

func (u *UserBalanceRepositoryMem) Upset(row model.UserBalance) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.userBalance[row.UserID] = row
	return nil
}
