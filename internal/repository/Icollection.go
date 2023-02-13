package repository

type ICollection interface {
	GetUserRepository() IUserRepository
	GetUserOrderRepository() IUserOrderRepository
	GetUserBalanceRepository() IUserBalanceRepository
	GetUserBalanceWithdrawRepository() IUserBalanceWithdrawRepository
}
