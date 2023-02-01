package handlers

import (
	"log"
	"net/http"

	"github.com/volkoviimagnit/gofermart/internal/handlers/response"
)

type AbstractHandler struct {
	httpMethod  string
	httpPattern string
}

func NewAbstractHandler(httpMethod string, httpPattern string) *AbstractHandler {
	return &AbstractHandler{
		httpMethod:  httpMethod,
		httpPattern: httpPattern,
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
	rw.Header().Set("Content-Type", "application/json")
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
