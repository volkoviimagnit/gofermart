package client

import "time"

type IAccrualClient interface {
	GetOrderStatus(orderNumber string) (IAccrualOrderStatus, IError)
	GetDefaultRetryAfterSeconds() time.Duration
}

type IError interface {
	NeedRetry() bool
	RetryAfterSeconds() time.Duration
}
