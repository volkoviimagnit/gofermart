package service

import "github.com/volkoviimagnit/gofermart/internal/transport"

type UserBalanceRecalculateMessage struct {
	ByOrderNumber string
}

func (u UserBalanceRecalculateMessage) GetQueueName() string {
	return transport.UserBalanceRecalculate
}
