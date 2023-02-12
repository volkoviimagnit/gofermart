package service

import "github.com/volkoviimagnit/gofermart/internal/client/response"

type IUserOrderService interface {
	AddOrder(userID string, orderNumber string) error
	Update(orderNumber string, status response.OrderStatus, accrual *float64) error
}
