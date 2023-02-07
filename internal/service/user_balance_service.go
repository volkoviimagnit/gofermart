package service

import (
	"fmt"
	"time"

	"github.com/volkoviimagnit/gofermart/internal/helpers"
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

// AddUserWithdraw TODO: расчеты нужно делать в рамках одной транзакции
func (u *UserBalanceService) AddUserWithdraw(userId string, orderNumber string, sum float64) error {
	if !helpers.ValidateOrderNumber(orderNumber) {
		return NewIncorrectOrderNumberError(orderNumber)
	}

	userBalance, errFindBalancing := u.GetUserBalance(userId)
	if errFindBalancing != nil {
		return NewBalanceNotFoundError(errFindBalancing.Error())
	}

	current := userBalance.GetCurrent()
	if current < sum {
		return NewNotEnoughFundsError(fmt.Sprintf("%f < %f", current, sum))
	}

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

	newBalance, errBalancing := u.SetUserBalance(userId, 0, 0)
	if errBalancing != nil {
		return nil, errBalancing
	}
	return newBalance, nil
}

func (u *UserBalanceService) SetUserBalance(userId string, current float64, withdrawn float64) (IUserBalance, error) {
	newBalance := model.NewUserBalance(userId, current, withdrawn)
	errInserting := u.userBalanceRepository.Upset(*newBalance)
	if errInserting != nil {
		return nil, errInserting
	}
	return newBalance, nil
}

type BalanceNotFoundError struct {
	err string
}

type IncorrectOrderNumberError struct {
	err string
}

type NotEnoughFundsError struct {
	err string
}

func NewBalanceNotFoundError(err string) *BalanceNotFoundError {
	return &BalanceNotFoundError{err: err}
}

func (e *BalanceNotFoundError) Error() string {
	return "баланс пользователя не найден - " + e.err
}

func NewIncorrectOrderNumberError(err string) *IncorrectOrderNumberError {
	return &IncorrectOrderNumberError{err: err}
}

func (e *IncorrectOrderNumberError) Error() string {
	return "неверный номер заказа - " + e.err
}

func NewNotEnoughFundsError(err string) *NotEnoughFundsError {
	return &NotEnoughFundsError{err: err}
}

func (e *NotEnoughFundsError) Error() string {
	return "на счету недостаточно средств - " + e.err
}
