package service

import (
	"github.com/volkoviimagnit/gofermart/internal/client"
	"github.com/volkoviimagnit/gofermart/internal/transport"
)

type OrderAccrualRequestMessage struct {
	OrderNumber        string
	AccrualOrderStatus client.IAccrualOrderStatus
}

func (o *OrderAccrualRequestMessage) GetQueueName() string {
	return transport.OrderAccrualQueueName
}
