package transport

type IMessenger interface {
	AddConsumer(consumer IConsumer)
	Dispatch(message IMessage)
	Consume(idx int, queueName string)
}
