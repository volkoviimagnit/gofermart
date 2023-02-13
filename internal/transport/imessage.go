package transport

type IMessage interface {
	GetQueueName() string
}
