package service

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/volkoviimagnit/gofermart/internal/transport"
)

type UserBalanceRecalculateConsumer struct {
	userBalanceService IUserBalanceService
}

func NewUserBalanceRecalculateConsumer(userBalanceService IUserBalanceService) transport.IConsumer {
	return &UserBalanceRecalculateConsumer{
		userBalanceService: userBalanceService,
	}
}

func (u *UserBalanceRecalculateConsumer) Execute(message transport.IMessage) error {
	logrus.Debug("UserBalanceRecalculateConsumer Execute")
	m, ok := message.(*UserBalanceRecalculateMessage)
	if !ok {
		return errors.New("UserBalanceRecalculateConsumer - конзюмер не поддерживает данный тип сообщений")
	}

	orderId := m.ByOrderNumber
	logrus.Debugf("UserBalanceRecalculateConsumer orderId:%s", orderId)
	err := u.userBalanceService.RecalculateByOrderNumber(orderId)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserBalanceRecalculateConsumer) GetQueueName() string {
	return transport.UserBalanceRecalculate
}
