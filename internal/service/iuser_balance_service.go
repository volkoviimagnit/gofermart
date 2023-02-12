package service

type IUserBalanceService interface {
	AddUserWithdraw(userID string, orderNumber string, sum float64) error
	GetUserBalance(userID string) (IUserBalance, error)
	SetUserBalance(userID string, current float64, withdrawn float64) (IUserBalance, error)
	RecalculateByOrderNumber(orderNumber string) error
	RecalculateByUserID(userID string) error
}
