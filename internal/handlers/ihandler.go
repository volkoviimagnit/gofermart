package handlers

import "net/http"

type IHandler interface {
	GetMethod() string
	GetPattern() string
	ServeHTTP(http.ResponseWriter, *http.Request)
}
