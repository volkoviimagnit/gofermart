package model

import "time"

type UserBalanceWithdraw struct {
	UserID      string
	OrderNumber string
	Sum         float64
	ProcessedAt time.Time
}

func (m *UserBalanceWithdraw) GetUserID() string {
	return m.UserID
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
