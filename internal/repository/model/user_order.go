package model

import "time"

type UserOrder struct {
	userId     string
	number     string
	status     string
	accrual    *float64
	uploadedAt time.Time
}

func (u *UserOrder) UserId() string {
	return u.userId
}

func (u *UserOrder) SetUserId(userId string) {
	u.userId = userId
}

func (u *UserOrder) Number() string {
	return u.number
}

func (u *UserOrder) SetNumber(number string) {
	u.number = number
}

func (u *UserOrder) Status() string {
	return u.status
}

func (u *UserOrder) SetStatus(status string) {
	u.status = status
}

func (u *UserOrder) Accrual() *float64 {
	return u.accrual
}

func (u *UserOrder) SetAccrual(accrual *float64) {
	u.accrual = accrual
}

func (u *UserOrder) UploadedAt() time.Time {
	return u.uploadedAt
}

func (u *UserOrder) SetUploadedAt(uploadedAt time.Time) {
	u.uploadedAt = uploadedAt
}
