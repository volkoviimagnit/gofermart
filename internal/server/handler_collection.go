package server

type HandlerCollection struct {
	handlers []IHttpHandler
}

func NewHandlerCollection() *HandlerCollection {
	return &HandlerCollection{
		handlers: make([]IHttpHandler, 0),
	}
}

func (c *HandlerCollection) AddHandler(handler IHttpHandler) ICollection {
	c.handlers = append(c.handlers, handler)
	return c
}

func (c *HandlerCollection) GetHandlers() []IHttpHandler {
	return c.handlers
}
