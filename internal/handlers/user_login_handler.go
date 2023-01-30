package handlers

import (
	"net/http"
)

type UserLoginHandler struct {
	parent *AbstractHandler
}

func (u UserLoginHandler) GetMethod() string {
	return u.parent.GetMethod()
}

func (u UserLoginHandler) GetPattern() string {
	return u.parent.GetPattern()
}

func (u UserLoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// http.StatusOK
	// http.StatusBadRequest
	// http.StatusUnauthorized
	// http.StatusInternalServerError

	u.parent.RenderResponse(writer, http.StatusOK, []byte("UserLoginHandler"))
}

func NewUserLoginHandler() IHandler {
	return &UserLoginHandler{
		parent: NewAbstractHandler(http.MethodPost, "/api/user/login"),
	}
}
