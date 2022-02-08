package chttp

import (
	"github.com/go-chi/chi"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, *Context)

func NewRouter() *Router {
	return &Router{
		Mux: chi.NewRouter(),
	}
}

type Router struct {
	*chi.Mux
	ehNotFound         HandlerFunc
	ehMethodNotAllowed HandlerFunc
}

func (r *Router) Route(pattern string, fn func(r *Router)) *Router {
	subRouter := NewRouter()
	if r.ehMethodNotAllowed != nil {
		subRouter.MethodNotAllowed(r.ehMethodNotAllowed)
	}
	if r.ehNotFound != nil {
		subRouter.NotFound(r.ehNotFound)
	}

	if fn != nil {
		fn(subRouter)
	}
	r.Mux.Mount(pattern, subRouter)

	return subRouter
}

func (r *Router) Mount(pattern string, h HandlerFunc) {
	r.Mux.Mount(pattern, r.wrapHandlerFunc(h))
}

func (r *Router) NotFound(h HandlerFunc) {
	r.ehNotFound = h
	r.Mux.NotFound(r.wrapHandlerFunc(h))
}

func (r *Router) MethodNotAllowed(h HandlerFunc) {
	r.ehMethodNotAllowed = h
	r.Mux.MethodNotAllowed(r.wrapHandlerFunc(h))
}

func (r *Router) HandleFunc(pattern string, h HandlerFunc) {
	r.Mux.HandleFunc(pattern, r.wrapHandlerFunc(h))
}

func (r *Router) Get(pattern string, h HandlerFunc) {
	r.Mux.Get(pattern, r.wrapHandlerFunc(h))
}

func (r *Router) Post(pattern string, h HandlerFunc) {
	r.Mux.Post(pattern, r.wrapHandlerFunc(h))
}

func (r *Router) Delete(pattern string, h HandlerFunc) {
	r.Mux.Delete(pattern, r.wrapHandlerFunc(h))
}

func (r *Router) Put(pattern string, h HandlerFunc) {
	r.Mux.Put(pattern, r.wrapHandlerFunc(h))
}

func (r *Router) Patch(pattern string, h HandlerFunc) {
	r.Mux.Patch(pattern, r.wrapHandlerFunc(h))
}

func (r *Router) wrapHandlerFunc(h HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := NewContext(res, req)

		h(res, req, ctx)
	}
}

func (r *Router) WrapHandler(handler http.Handler) HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request, ctx *Context) {
		handler.ServeHTTP(res, req)
	}
}
