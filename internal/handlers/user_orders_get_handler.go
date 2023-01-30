package handlers

import "net/http"

type UserOrdersGETHandler struct {
	parent *AbstractHandler
}

func NewUserOrdersGETHandler() *UserOrdersGETHandler {
	return &UserOrdersGETHandler{
		parent: NewAbstractHandler(http.MethodGet, "/api/user/orders"),
	}
}

func (h *UserOrdersGETHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserOrdersGETHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserOrdersGETHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusNoContent
	// http.StatusUnauthorized
	// http.StatusInternalServerError

	h.parent.RenderResponse(rw, http.StatusOK, []byte("UserOrdersGETHandler"))
}
