package client

import "github.com/volkoviimagnit/gofermart/internal/client/response"

type IAccrualOrderStatus interface {
	GetNumber() string
	GetStatus() response.OrderStatus
	GetAccrual() *float64

	IsTerminal() bool
}
