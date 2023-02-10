package transport

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type MessengerMem struct {
	queue          map[string][]IMessage
	queueConsumers []IConsumer
	mutex          *sync.RWMutex
}

const OrderAccrualQueueName = "order_accrual_request"
const UserBalanceRecalculate = "user_balance_recalculate"

func NewMessengerMem() IMessenger {
	return &MessengerMem{
		queue:          make(map[string][]IMessage, 5),
		queueConsumers: make([]IConsumer, 0),
		mutex:          &sync.RWMutex{},
	}
}

func (m *MessengerMem) AddConsumer(consumer IConsumer) {
	m.queueConsumers = append(m.queueConsumers, consumer)
}

func (m *MessengerMem) Dispatch(message IMessage) {
	m.queue[message.GetQueueName()] = append(m.queue[message.GetQueueName()], message)
	/*
		go func() {
			errRunning := message.Run()
			if errRunning != nil {
				logrus.Fatalf("goroutine error %s", errRunning)
				return
			}
			return
		}()
	*/
	//m.Signal(message.GetQueueName())
}

// Consume TODO: узнать как это правильно реализовать через каналы?
func (m *MessengerMem) Consume(idx int, queueName string) {
	for _, consumer := range m.queueConsumers {
		c := consumer
		go func() {
			for {
				<-time.After(1 * time.Second)
				qName := c.GetQueueName()
				//logrus.Debugf("Опрос очереди %d - %s", consumerId, qName)
				firstMessage := m.QueuePop(qName)
				if firstMessage != nil {
					errRunning := c.Execute(firstMessage)
					if errRunning != nil {
						logrus.Fatalf("goroutine error %s", errRunning)
					}
				}
			}
		}()
	}
}

func (m *MessengerMem) QueuePop(queueName string) IMessage {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	queue := m.queue[queueName]
	if len(queue) > 0 {
		firstMessage := queue[0]
		if len(queue) > 1 {
			m.queue[queueName] = m.queue[queueName][1:]
		} else {
			m.queue[queueName] = m.queue[queueName][:0]
		}
		return firstMessage
	}
	return nil
}
