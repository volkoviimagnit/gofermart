package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type UserOrderRepositoryMem struct {
	userOrders map[string][]model.UserOrder
}

func NewUserOrderRepositoryMem() IUserOrderRepository {
	return &UserOrderRepositoryMem{userOrders: make(map[string][]model.UserOrder, 0)}
}

func (r *UserOrderRepositoryMem) Insert(row model.UserOrder) error {
	if _, isExist := r.userOrders[row.UserId()]; !isExist {
		r.userOrders[row.UserId()] = make([]model.UserOrder, 0)
	}
	r.userOrders[row.UserId()] = append(r.userOrders[row.UserId()], row)
	return nil
}

func (r *UserOrderRepositoryMem) FindByUserId(userId string) ([]model.UserOrder, error) {
	if _, isExist := r.userOrders[userId]; !isExist {
		return make([]model.UserOrder, 0), nil
	}
	return r.userOrders[userId], nil
}
