package service

import "github.com/volkoviimagnit/gofermart/internal/client/response"

type IUserOrderService interface {
	AddOrder(userId string, orderId string) error
	Update(orderId string, status response.OrderStatus, accrual *float64) error
}
