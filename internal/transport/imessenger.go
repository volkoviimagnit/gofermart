package transport

type IMessenger interface {
	AddConsumer(consumer IConsumer)
	Dispatch(message IMessage)
	Consume()
}
