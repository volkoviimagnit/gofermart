package service

type IUserBalanceService interface {
	AddUserWithdraw(userId string, orderNumber string, sum float64) error
	GetUserBalance(userId string) (IUserBalance, error)
}
