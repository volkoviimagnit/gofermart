package service

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/volkoviimagnit/gofermart/internal/client"
	"github.com/volkoviimagnit/gofermart/internal/client/response"
	"github.com/volkoviimagnit/gofermart/internal/repository"
	"github.com/volkoviimagnit/gofermart/internal/repository/model"
	"github.com/volkoviimagnit/gofermart/internal/transport"
)

type UserOrderService struct {
	accrualClient                 client.IAccrualClient
	messenger                     transport.IMessenger
	userOrderRepository           repository.IUserOrderRepository
	userBalanceRepository         repository.IUserBalanceRepository
	userBalanceWithdrawRepository repository.IUserBalanceWithdrawRepository
}

func (u *UserOrderService) Update(orderNumber string, status response.OrderStatus, accrual *float64) error {
	logrus.Debugf("UserOrderService.Update")
	oldOrder, errRepository := u.userOrderRepository.FindOneByNumber(orderNumber)
	if errRepository != nil {
		return &RepositoryError{err: errRepository}
	}
	userOrderStatus := u.generateUserOrderStatus(status)
	oldOrder.Status = userOrderStatus
	oldOrder.UploadedAt = time.Now()
	oldOrder.Accrual = accrual
	errUpdating := u.userOrderRepository.Update(*oldOrder)
	if errUpdating != nil {
		return errUpdating
	}
	return nil
}

func (u *UserOrderService) generateUserOrderStatus(accrualStatus response.OrderStatus) model.UserOrderStatus {
	switch accrualStatus {
	case response.AccrualStatusRegistered:
		return model.UserOrderStatusNew
	case response.AccrualStatusProcessing:
		return model.UserOrderStatusProcessing
	case response.AccrualStatusProcessed:
		return model.UserOrderStatusProcessed
	case response.AccrualStatusInvalid:
		return model.UserOrderStatusInvalid
	default:
		return model.UserOrderStatusInvalid
	}
}

func (u *UserOrderService) AddOrder(userID string, orderNumber string) error {
	oldOrder, errRepository := u.userOrderRepository.FindOneByNumber(orderNumber)
	if errRepository != nil {
		return &RepositoryError{err: errRepository}
	}
	if oldOrder != nil {
		isOwnOrder := oldOrder.UserID == userID
		if isOwnOrder {
			return &DuplicatedOwnOrderError{}
		} else {
			return &DuplicatedSomebodyElseOrderError{}
		}
	}

	accrual := 0.0
	m := model.UserOrder{
		UserID:     userID,
		Number:     orderNumber,
		Status:     model.UserOrderStatusNew,
		Accrual:    &accrual,
		UploadedAt: time.Now(),
	}

	errInserting := u.userOrderRepository.Insert(m)
	if errInserting != nil {
		return &RepositoryError{err: errInserting}
	}

	mess := OrderAccrualRequestMessage{
		OrderNumber: orderNumber,
	}
	u.messenger.Dispatch(&mess)
	return nil
}

func NewUserOrderService(
	accrualClient client.IAccrualClient,
	messenger transport.IMessenger,
	userOrderRepository repository.IUserOrderRepository,
	userBalanceRepository repository.IUserBalanceRepository,
	userBalanceWithdrawRepository repository.IUserBalanceWithdrawRepository,
) *UserOrderService {
	return &UserOrderService{
		accrualClient:                 accrualClient,
		messenger:                     messenger,
		userOrderRepository:           userOrderRepository,
		userBalanceRepository:         userBalanceRepository,
		userBalanceWithdrawRepository: userBalanceWithdrawRepository,
	}
}

type RepositoryError struct {
	err error
}

func (e *RepositoryError) Error() string {
	return "ошибка репозитория - " + e.err.Error()
}

type DuplicatedOwnOrderError struct {
}

func (e *DuplicatedOwnOrderError) Error() string {
	return "дубль собственного заказа"
}

type DuplicatedSomebodyElseOrderError struct {
}

func (e *DuplicatedSomebodyElseOrderError) Error() string {
	return "дубль собственного заказа"
}
