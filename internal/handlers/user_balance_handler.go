package handlers

import "net/http"

type UserBalanceHandler struct {
	parent *AbstractHandler
}

func NewUserBalanceHandler() *UserBalanceHandler {
	return &UserBalanceHandler{
		parent: NewAbstractHandler(http.MethodGet, "/api/user/balance", "application/json"),
	}
}

func (h *UserBalanceHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserBalanceHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserBalanceHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusUnauthorized
	// http.StatusInternalServerError

	h.parent.RenderResponse(rw, http.StatusOK, []byte("UserBalanceHandler"))
}
