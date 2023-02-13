package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
	"github.com/volkoviimagnit/gofermart/internal/security"
)

type AbstractHandler struct {
	httpMethod  string
	httpPattern string
	contentType string
	auth        security.IAuthenticator
}

func NewAbstractHandler(httpMethod string, httpPattern string, contentType string) *AbstractHandler {
	return &AbstractHandler{
		httpMethod:  httpMethod,
		httpPattern: httpPattern,
		contentType: contentType,
	}
}

func (h *AbstractHandler) Render(rw http.ResponseWriter, resp *response.Response) {
	rw.Header().Set("Content-Type", resp.GetContentType())
	rw.WriteHeader(resp.GetStatus())
	_, err := rw.Write(resp.GetBody())
	if err != nil {
		log.Fatal("rw.Write error in update")
	}
}

func (h *AbstractHandler) RenderResponse(rw http.ResponseWriter, statusCode int, body []byte) {
	rw.Header().Set("Content-Type", h.contentType)
	rw.WriteHeader(statusCode)
	_, err := rw.Write(body)
	if err != nil {
		log.Fatal("rw.Write error in update")
	}
}

func (h *AbstractHandler) GetMethod() string {
	return h.httpMethod
}

func (h *AbstractHandler) GetPattern() string {
	return h.httpPattern
}

func (h *AbstractHandler) RenderUnauthorized(rw http.ResponseWriter) {
	resp := response.NewResponse(h.contentType)
	resp.SetStatus(http.StatusUnauthorized).SetBody([]byte(""))
	h.Render(rw, resp)
}

func (h *AbstractHandler) RenderInternalServerError(rw http.ResponseWriter, err error) {
	resp := response.NewResponse(h.contentType)
	resp.SetStatus(http.StatusInternalServerError).SetBody([]byte(err.Error()))
	h.Render(rw, resp)
}

func (h *AbstractHandler) RenderError(rw http.ResponseWriter, status int, err error) {
	resp := response.NewResponse(h.contentType)
	resp.SetStatus(status).SetBody([]byte(err.Error()))
	h.Render(rw, resp)
}

func (h *AbstractHandler) AuthOrAbort(rw http.ResponseWriter, request *http.Request) security.IPassport {
	if h.auth == nil {
		h.RenderInternalServerError(rw, errors.New("способ аутентификации не задан"))
		return nil
	}

	passport, errAuth := h.auth.Authenticate(request)
	if errAuth != nil {
		h.RenderInternalServerError(rw, errAuth)
		return nil
	}
	if passport == nil {
		h.RenderUnauthorized(rw)
		return nil
	}
	return passport
}

func (h *AbstractHandler) SetAuthenticator(auth security.IAuthenticator) {
	h.auth = auth
}

func (h *AbstractHandler) RenderNoContent(rw http.ResponseWriter) {
	resp := response.NewResponse(h.contentType)
	resp.SetStatus(http.StatusNoContent).SetBody([]byte(""))
	h.Render(rw, resp)
}
