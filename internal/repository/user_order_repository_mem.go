package repository

import (
	"sync"

	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserOrderRepositoryMem struct {
	userOrders map[string]map[string]model.UserOrder
	mutex      *sync.RWMutex
}

func NewUserOrderRepositoryMem() IUserOrderRepository {
	return &UserOrderRepositoryMem{
		userOrders: make(map[string]map[string]model.UserOrder, 0),
		mutex:      &sync.RWMutex{},
	}
}

func (r *UserOrderRepositoryMem) SumAccrualByUserID(userID string) (float64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	sum := 0.0
	for _, userOrder := range r.userOrders[userID] {
		if userOrder.Accrual == nil {
			continue
		}
		if !userOrder.Status.IsSuitable() {
			continue
		}
		sum += *userOrder.Accrual
	}
	return sum, nil
}

func (r *UserOrderRepositoryMem) Insert(row model.UserOrder) error {
	return r.Upsert(row)
}

func (r *UserOrderRepositoryMem) Update(row model.UserOrder) error {
	return r.Upsert(row)
}

func (r *UserOrderRepositoryMem) Upsert(row model.UserOrder) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, isExist := r.userOrders[row.UserID]; !isExist {
		r.userOrders[row.UserID] = make(map[string]model.UserOrder, 0)
	}
	r.userOrders[row.UserID][row.Number] = row
	return nil
}

func (r *UserOrderRepositoryMem) FindByUserID(userID string) ([]model.UserOrder, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if _, isExist := r.userOrders[userID]; !isExist {
		return make([]model.UserOrder, 0), nil
	}
	if len(r.userOrders[userID]) == 0 {
		return make([]model.UserOrder, 0), nil
	}
	userOrders := make([]model.UserOrder, 0, len(r.userOrders[userID]))
	for _, userOrder := range r.userOrders[userID] {
		userOrders = append(userOrders, userOrder)
	}
	return userOrders, nil
}

func (r *UserOrderRepositoryMem) FindOneByNumber(number string) (*model.UserOrder, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, orders := range r.userOrders {
		for _, order := range orders {
			if order.Number == number {
				return &order, nil
			}
		}
	}
	return nil, nil
}
