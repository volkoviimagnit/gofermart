package handlers

import (
	"io"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/request"
	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/security"
)

type UserOrdersPOSTHandler struct {
	parent *AbstractHandler
	auth   security.IAuthenticator
}

func NewUserOrderPOSTHandler(auth security.IAuthenticator) *UserOrdersPOSTHandler {
	return &UserOrdersPOSTHandler{
		parent: NewAbstractHandler(http.MethodPost, "/api/user/orders"),
		auth:   auth,
	}
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

	resp := response.NewResponse("text/plain")

	passport, errAuth := h.auth.Authenticate(request)
	if errAuth != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(errAuth.Error()))
		h.parent.Render(rw, resp)
		return
	}
	if passport == nil {
		resp.SetStatus(http.StatusUnauthorized).SetBody([]byte("пользователь не авторизован"))
		h.parent.Render(rw, resp)
		return
	}

	dto, errBody := h.extractRequestDTO(request)
	if errBody != nil {
		resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(errBody.Error()))
		h.parent.Render(rw, resp)
		return
	}

	errValidation := dto.Validate()
	if errValidation != nil {
		resp.SetStatus(http.StatusBadRequest).SetBody([]byte(errValidation.Error()))
		h.parent.Render(rw, resp)
		return
	}

	resp.SetStatus(http.StatusAccepted).SetBody([]byte(dto.GetNumber()))
	h.parent.Render(rw, resp)
	return
}

func (h *UserOrdersPOSTHandler) extractRequestDTO(r *http.Request) (*request.UserOrdersPOSTDTO, error) {
	body, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		return nil, errBody
	}
	requestDTO := request.NewUserOrdersPOSTDTO(string(body))
	return requestDTO, nil
}
