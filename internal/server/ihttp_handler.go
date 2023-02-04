package server

import "net/http"

type IHttpHandler interface {
	GetMethod() string
	GetPattern() string
	ServeHTTP(http.ResponseWriter, *http.Request)
}
