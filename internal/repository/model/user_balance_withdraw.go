package model

import "time"

type UserBalanceWithdraw struct {
	UserId      string
	OrderNumber string
	Sum         float64
	ProcessedAt time.Time
}

func (m *UserBalanceWithdraw) GetUserId() string {
	return m.UserId
}

func (m *UserBalanceWithdraw) GetOrderNumber() string {
	return m.OrderNumber
}

func (m *UserBalanceWithdraw) GetSum() float64 {
	return m.Sum
}

func (m *UserBalanceWithdraw) GetProcessedAt() time.Time {
	return m.ProcessedAt
}
