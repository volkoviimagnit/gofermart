package handlers

import "net/http"

type OrderNumberHandler struct {
	parent *AbstractHandler
}

func NewOrderNumberHandler() *OrderNumberHandler {
	return &OrderNumberHandler{
		parent: NewAbstractHandler(http.MethodGet, "/api/orders/{number:([0-9]+)(?i)}"),
	}
}

func (h *OrderNumberHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *OrderNumberHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *OrderNumberHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusTooManyRequests
	// http.StatusInternalServerError

	h.parent.RenderResponse(rw, http.StatusOK, []byte("OrderNumberHandler"))
}
