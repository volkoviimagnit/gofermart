package server

import "net/http"

type IRouter interface {
	Configure() error
	GetHandler() http.Handler
}
