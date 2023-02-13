package transport

type IConsumer interface {
	Execute(message IMessage) error
	GetQueueName() string
}
