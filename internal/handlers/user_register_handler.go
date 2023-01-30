package handlers

import (
	"net/http"
)

type UserRegisterHandler struct {
	parent *AbstractHandler
}

func NewUserRegisterHandler() IHandler {
	return &UserRegisterHandler{
		parent: NewAbstractHandler(http.MethodPost, "/api/user/register"),
	}
}

func (h *UserRegisterHandler) GetMethod() string {
	return h.parent.GetMethod()
}

func (h *UserRegisterHandler) GetPattern() string {
	return h.parent.GetPattern()
}

func (h *UserRegisterHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusBadRequest
	// http.StatusConflict
	// http.StatusInternalServerError

	h.parent.RenderResponse(rw, http.StatusOK, []byte("UserRegisterHandler"))
}
