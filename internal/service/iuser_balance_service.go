package service

type IUserBalanceService interface {
	AddUserWithdraw(userId string, orderNumber string, sum float64) error
	GetUserBalance(userId string) (IUserBalance, error)
	SetUserBalance(userId string, current float64, withdrawn float64) (IUserBalance, error)
}
