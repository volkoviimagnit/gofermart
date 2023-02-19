package model

import "time"

type UserOrder struct {
	UserID     string
	Number     string
	Status     UserOrderStatus
	Accrual    *float64
	UploadedAt time.Time
}
