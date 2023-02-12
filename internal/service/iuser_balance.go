package service

type IUserBalance interface {
	GetCurrent() float64
	GetWithdrawn() float64
	SetCurrent(current float64)
	SetWithdrawn(Withdrawn float64)
}
