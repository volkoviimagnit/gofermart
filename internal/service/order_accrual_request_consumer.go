package service

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/volkoviimagnit/gofermart/internal/client"
	"github.com/volkoviimagnit/gofermart/internal/transport"
)

type OrderAccrualRequestConsumer struct {
	messenger        transport.IMessenger
	httpClient       client.IAccrualClient
	userOrderService IUserOrderService
}

func NewOrderAccrualConsumer(
	messenger transport.IMessenger,
	httpClient client.IAccrualClient,
	userOrderService IUserOrderService,
) transport.IConsumer {
	return &OrderAccrualRequestConsumer{
		messenger:        messenger,
		httpClient:       httpClient,
		userOrderService: userOrderService,
	}
}

func (o *OrderAccrualRequestConsumer) GetQueueName() string {
	return transport.OrderAccrualQueueName
}

func (o *OrderAccrualRequestConsumer) Execute(message transport.IMessage) error {
	logrus.Debug("OrderAccrualRequestConsumer Execute")
	m, ok := message.(*OrderAccrualRequestMessage)
	if !ok {
		return errors.New("конзюмер не поддерживает данный тип сообщений")
	}
	orderNumber := m.OrderNumber

	logrus.Debugf("Запрос к accrual /%s", orderNumber)
	orderStatus, errRequesting := o.httpClient.GetOrderStatus(orderNumber)
	if errRequesting != nil {
		logrus.Debugf("Получена ошибка %+v", errRequesting)
		if errRequesting.NeedRetry() {
			o.retryRequest(orderNumber, errRequesting.RetryAfterSeconds())
		} else {
			o.cancelUserOrder(orderNumber)
		}
		return nil
	}

	if !orderStatus.IsTerminal() {
		o.retryRequest(orderNumber, o.httpClient.GetDefaultRetryAfterSeconds())
		return nil
	}

	errUpdating := o.userOrderService.Update(orderNumber, orderStatus.GetStatus(), orderStatus.GetAccrual())
	if errUpdating != nil {
		logrus.Errorf("не удалось обновить статус заказа - %s", errUpdating.Error())
		o.retryRequest(orderNumber, o.httpClient.GetDefaultRetryAfterSeconds())
		return nil
	}

	// пересчет запускаем не по владельцу, а отталкиваясь от заказа
	o.messenger.Dispatch(&UserBalanceRecalculateMessage{
		ByOrderNumber: orderNumber,
	})

	return nil

}

// retryRequest TODO: уточнить в ТЗ что с retry-policy?
func (o *OrderAccrualRequestConsumer) retryRequest(orderNumber string, retryAfter time.Duration) {
	logrus.Debugf("отложенное выполнение %s", orderNumber)
	time.AfterFunc(retryAfter*time.Second, func() {
		o.messenger.Dispatch(&OrderAccrualRequestMessage{
			OrderNumber: orderNumber,
		})
	})

}

func (o *OrderAccrualRequestConsumer) cancelUserOrder(orderNumber string) {
	logrus.Warningf("Заказа не существует. Отменяем заказ.")

	o.messenger.Dispatch(&OrderAccrualRequestMessage{
		OrderNumber: orderNumber,
	})
}
