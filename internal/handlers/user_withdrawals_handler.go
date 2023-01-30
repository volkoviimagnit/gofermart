package handlers

import "net/http"

type UserWithdrawalsHandler struct {
	parent *AbstractHandler
}

func NewUserWithdrawalsHandler() *UserWithdrawalsHandler {
	return &UserWithdrawalsHandler{
		parent: NewAbstractHandler(http.MethodGet, "/api/user/withdrawals"),
	}
}

func (h *UserWithdrawalsHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserWithdrawalsHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserWithdrawalsHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusNoContent
	// http.StatusUnauthorized
	// http.StatusInternalServerError

	h.parent.RenderResponse(rw, http.StatusOK, []byte("UserWithdrawalsHandler"))
}
