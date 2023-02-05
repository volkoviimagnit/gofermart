package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterChi struct {
	isDebugMode       bool
	handlerCollection ICollection
	mux               *chi.Mux
}

func NewRouterChi(handlerCollection ICollection, isDebugMode bool) IRouter {
	return &RouterChi{
		mux:               chi.NewRouter(),
		isDebugMode:       isDebugMode,
		handlerCollection: handlerCollection,
	}

}
func (r *RouterChi) Configure() error {
	r.mux.Use(middleware.RequestID)
	r.mux.Use(middleware.RealIP)
	if r.isDebugMode {
		r.mux.Use(middleware.Logger)
	}
	r.mux.Use(middleware.Recoverer)
	r.mux.Use(middleware.StripSlashes)
	r.mux.Use(middleware.Timeout(60 * time.Second))

	for _, handler := range r.handlerCollection.GetHandlers() {
		r.mux.Method(handler.GetMethod(), handler.GetPattern(), handler)
	}

	return nil
}

func (r *RouterChi) GetHandler() http.Handler {
	return r.mux
}
