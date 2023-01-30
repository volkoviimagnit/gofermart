package handlers

import "net/http"

type UserOrdersPOSTHandler struct {
	parent *AbstractHandler
}

func NewUserOrderPOSTHandler() *UserOrdersPOSTHandler {
	return &UserOrdersPOSTHandler{parent: NewAbstractHandler(http.MethodPost, "/api/user/orders")}
}

func (h *UserOrdersPOSTHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserOrdersPOSTHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserOrdersPOSTHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	//http.StatusAccepted
	//http.StatusBadRequest
	//http.StatusUnauthorized
	//http.StatusConflict
	//http.StatusUnprocessableEntity
	// http.StatusInternalServerError

	h.parent.RenderResponse(rw, http.StatusOK, []byte("UserOrdersPOSTHandler"))
}
