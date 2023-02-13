package server

type ICollection interface {
	AddHandler(handler IHttpHandler) ICollection
	GetHandlers() []IHttpHandler
}
