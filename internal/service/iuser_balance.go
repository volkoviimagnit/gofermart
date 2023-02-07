package service

type IUserBalance interface {
	GetCurrent() float64
	GetWithdrawn() float64
}
