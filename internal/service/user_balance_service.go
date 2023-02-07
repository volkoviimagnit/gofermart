package service

import (
	"time"

	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
)

type UserBalanceService struct {
	userBalanceRepository         repository.IUserBalanceRepository
	userBalanceWithdrawRepository repository.IUserBalanceWithdrawRepository
}

func NewUserBalanceService(
	ubRepository repository.IUserBalanceRepository,
	ubwRepository repository.IUserBalanceWithdrawRepository) IUserBalanceService {

	return &UserBalanceService{
		userBalanceRepository:         ubRepository,
		userBalanceWithdrawRepository: ubwRepository}
}

// AddUserWithdraw TODO расчеты нужно делать в рамках одной транзакции
func (u *UserBalanceService) AddUserWithdraw(userId string, orderNumber string, sum float64) error {

	userBalanceWithdrawModel := model.UserBalanceWithdraw{
		UserId:      userId,
		OrderNumber: orderNumber,
		Sum:         sum,
		ProcessedAt: time.Now(),
	}
	errWithdrawInserting := u.userBalanceWithdrawRepository.Insert(userBalanceWithdrawModel)
	return errWithdrawInserting
}

func (u *UserBalanceService) GetUserBalance(userId string) (IUserBalance, error) {
	row, err := u.userBalanceRepository.FinOneByUserId(userId)
	if err != nil {
		return nil, err
	}
	if row != nil {
		return row, nil
	}

	newBalance := model.NewUserBalance(userId, 0, 0)
	errInserting := u.userBalanceRepository.Insert(*newBalance)
	if errInserting != nil {
		return nil, errInserting
	}
	return newBalance, nil
}
