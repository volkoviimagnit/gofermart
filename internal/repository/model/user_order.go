package model

import "time"

type UserOrder struct {
	userID     string
	number     string
	status     UserOrderStatus
	accrual    *float64
	uploadedAt time.Time
}

func (u *UserOrder) GetUserID() string {
	return u.userID
}

func (u *UserOrder) SetUserID(userID string) {
	u.userID = userID
}

func (u *UserOrder) GetNumber() string {
	return u.number
}

func (u *UserOrder) SetNumber(number string) {
	u.number = number
}

func (u *UserOrder) GetStatus() UserOrderStatus {
	return u.status
}

func (u *UserOrder) SetStatus(status UserOrderStatus) {
	u.status = status
}

func (u *UserOrder) GetAccrual() *float64 {
	return u.accrual
}

func (u *UserOrder) SetAccrual(accrual *float64) {
	u.accrual = accrual
}

func (u *UserOrder) GetUploadedAt() time.Time {
	return u.uploadedAt
}

func (u *UserOrder) SetUploadedAt(uploadedAt time.Time) {
	u.uploadedAt = uploadedAt
}
