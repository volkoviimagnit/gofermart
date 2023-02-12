package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type UserOrderRepositoryMem struct {
	userOrders map[string]map[string]model.UserOrder
}

func (r *UserOrderRepositoryMem) SumAccrualByUserID(userID string) (float64, error) {
	sum := 0.0
	for _, userOrder := range r.userOrders[userID] {
		if userOrder.GetAccrual() == nil {
			continue
		}
		if !userOrder.GetStatus().IsSuitable() {
			continue
		}
		sum += *userOrder.GetAccrual()
	}
	return sum, nil
}

func NewUserOrderRepositoryMem() IUserOrderRepository {
	return &UserOrderRepositoryMem{userOrders: make(map[string]map[string]model.UserOrder, 0)}
}

func (r *UserOrderRepositoryMem) Insert(row model.UserOrder) error {
	if _, isExist := r.userOrders[row.GetUserID()]; !isExist {
		r.userOrders[row.GetUserID()] = make(map[string]model.UserOrder, 0)
	}
	r.userOrders[row.GetUserID()][row.GetNumber()] = row
	return nil
}

func (r *UserOrderRepositoryMem) Update(row model.UserOrder) error {
	if _, isExist := r.userOrders[row.GetUserID()]; !isExist {
		r.userOrders[row.GetUserID()] = make(map[string]model.UserOrder, 0)
	}
	r.userOrders[row.GetUserID()][row.GetNumber()] = row
	return nil
}

func (r *UserOrderRepositoryMem) FindByUserID(userID string) ([]model.UserOrder, error) {
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
	for _, orders := range r.userOrders {
		for _, order := range orders {
			if order.GetNumber() == number {
				return &order, nil
			}
		}
	}
	return nil, nil
}
