package handlers

import "net/http"

type UserBalanceWithdrawHandler struct {
	parent *AbstractHandler
}

func NewUserBalanceWithdrawHandler() *UserBalanceWithdrawHandler {
	return &UserBalanceWithdrawHandler{
		parent: NewAbstractHandler(http.MethodPost, "/api/user/balance/withdraw", "application/json"),
	}
}

func (h *UserBalanceWithdrawHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserBalanceWithdrawHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserBalanceWithdrawHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusUnauthorized
	// http.StatusPaymentRequired
	// http.StatusUnprocessableEntity
	// http.StatusInternalServerError

	h.parent.RenderResponse(rw, http.StatusOK, []byte("UserBalanceWithdrawHandler"))
}
