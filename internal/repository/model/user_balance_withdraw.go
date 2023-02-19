package model

import "time"

type UserBalanceWithdraw struct {
	UserID      string
	OrderNumber string
	Sum         float64
	ProcessedAt time.Time
}
