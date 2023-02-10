package repository

import "github.com/volkoviimagnit/gofermart/internal/repository/model"

type UserOrderRepositoryMem struct {
	userOrders map[string]map[string]model.UserOrder
}

func (r *UserOrderRepositoryMem) SumAccrualByUserId(userId string) float64 {
	sum := 0.0
	for _, userOrder := range r.userOrders[userId] {
		if userOrder.Accrual() == nil {
			continue
		}
		if !userOrder.Status().IsSuitable() {
			continue
		}
		sum += *userOrder.Accrual()
	}
	return sum
}

func NewUserOrderRepositoryMem() IUserOrderRepository {
	return &UserOrderRepositoryMem{userOrders: make(map[string]map[string]model.UserOrder, 0)}
}

func (r *UserOrderRepositoryMem) Insert(row model.UserOrder) error {
	if _, isExist := r.userOrders[row.UserId()]; !isExist {
		r.userOrders[row.UserId()] = make(map[string]model.UserOrder, 0)
	}
	r.userOrders[row.UserId()][row.Number()] = row
	return nil
}

func (r *UserOrderRepositoryMem) Update(row model.UserOrder) error {
	if _, isExist := r.userOrders[row.UserId()]; !isExist {
		r.userOrders[row.UserId()] = make(map[string]model.UserOrder, 0)
	}
	r.userOrders[row.UserId()][row.Number()] = row
	return nil
}

func (r *UserOrderRepositoryMem) FindByUserId(userId string) ([]model.UserOrder, error) {
	if _, isExist := r.userOrders[userId]; !isExist {
		return make([]model.UserOrder, 0), nil
	}
	if len(r.userOrders[userId]) == 0 {
		return make([]model.UserOrder, 0), nil
	}
	userOrders := make([]model.UserOrder, 0, len(r.userOrders[userId]))
	for _, userOrder := range r.userOrders[userId] {
		userOrders = append(userOrders, userOrder)
	}
	return userOrders, nil
}

func (r *UserOrderRepositoryMem) FindOneByNumber(number string) (*model.UserOrder, error) {
	for _, orders := range r.userOrders {
		for _, order := range orders {
			if order.Number() == number {
				return &order, nil
			}
		}
	}
	return nil, nil
}

func (r *UserOrderRepositoryMem) IsExist(userId string, number string) (bool, error) {
	if _, isExist := r.userOrders[userId]; !isExist {
		return false, nil
	}
	if _, isExist := r.userOrders[userId][number]; !isExist {
		return false, nil
	}
	return true, nil
}
