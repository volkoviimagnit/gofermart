package service

import (
	"fmt"
	"time"

	"github.com/volkoviimagnit/gofermart/internal/helpers"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
	"github.com/volkoviimagnit/gofermart/internal/transport"
)

type UserBalanceService struct {
	userBalanceRepository         repository.IUserBalanceRepository
	userBalanceWithdrawRepository repository.IUserBalanceWithdrawRepository
	userOrderRepository           repository.IUserOrderRepository
	messenger                     transport.IMessenger
}

func NewUserBalanceService(
	ubRepository repository.IUserBalanceRepository,
	ubwRepository repository.IUserBalanceWithdrawRepository,
	userOrderRepository repository.IUserOrderRepository,
	messenger transport.IMessenger,
) IUserBalanceService {

	return &UserBalanceService{
		userBalanceRepository:         ubRepository,
		userBalanceWithdrawRepository: ubwRepository,
		userOrderRepository:           userOrderRepository,
		messenger:                     messenger,
	}
}

// AddUserWithdraw TODO: расчеты нужно делать в рамках одной транзакции
func (u *UserBalanceService) AddUserWithdraw(userID string, orderNumber string, sum float64) error {
	if !helpers.ValidateOrderNumber(orderNumber) {
		return NewIncorrectOrderNumberError(orderNumber)
	}

	userBalance, errFindBalancing := u.GetUserBalance(userID)
	if errFindBalancing != nil {
		return NewBalanceNotFoundError(errFindBalancing.Error())
	}

	current := userBalance.GetCurrent()
	if current < sum {
		return NewNotEnoughFundsError(fmt.Sprintf("%f < %f", current, sum))
	}

	userBalanceWithdrawModel := model.UserBalanceWithdraw{
		UserId:      userID,
		OrderNumber: orderNumber,
		Sum:         sum,
		ProcessedAt: time.Now(),
	}
	errWithdrawInserting := u.userBalanceWithdrawRepository.Insert(userBalanceWithdrawModel)
	if errWithdrawInserting != nil {
		return errWithdrawInserting
	}
	errRecalculating := u.RecalculateByUserID(userID)
	return errRecalculating
}

func (u *UserBalanceService) GetUserBalance(userID string) (IUserBalance, error) {
	row, err := u.userBalanceRepository.FinOneByUserID(userID)
	if err != nil {
		return nil, err
	}
	if row != nil {
		return row, nil
	}

	newBalance, errBalancing := u.SetUserBalance(userID, 0, 0)
	if errBalancing != nil {
		return nil, errBalancing
	}
	return newBalance, nil
}

func (u *UserBalanceService) SetUserBalance(userID string, current float64, withdrawn float64) (IUserBalance, error) {
	newBalance := model.NewUserBalance(userID, current, withdrawn)
	errInserting := u.userBalanceRepository.Upset(*newBalance)
	if errInserting != nil {
		return nil, errInserting
	}
	return newBalance, nil
}

func (u *UserBalanceService) RecalculateByOrderNumber(orderNumber string) error {
	userOrder, errOrder := u.userOrderRepository.FindOneByNumber(orderNumber)
	if errOrder != nil {
		return errOrder
	}
	userID := userOrder.UserId()
	return u.RecalculateByUserID(userID)
}

func (u *UserBalanceService) RecalculateByUserID(userID string) error {
	userBalance, errBalance := u.userBalanceRepository.FinOneByUserID(userID)
	if errBalance != nil {
		return errBalance
	}

	orderProcessedSum := u.userOrderRepository.SumAccrualByUserId(userID)
	orderWithdrawSum := u.userBalanceWithdrawRepository.SumWithdrawByUserId(userID)
	current := orderProcessedSum - orderWithdrawSum
	if current < 0 {
		current = 0
	}
	userBalance.SetCurrent(current)
	userBalance.SetWithdrawn(orderWithdrawSum)
	errBalanceUpdating := u.userBalanceRepository.Update(*userBalance)
	if errBalanceUpdating != nil {
		return errBalanceUpdating
	}
	return nil
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
