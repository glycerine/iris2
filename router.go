package iris2

import (
	"github.com/buaazp/fasthttprouter"
)

type HandlerFunc func(ctx *Context) error

type Router struct {
	path   string
	router *fasthttprouter.Router
}

func newRouter() *Router {
	return &Router{
		path:   "",
		router: fasthttprouter.New(),
	}
}

func (r *Router) Group(path string) *Router {
	return &Router{
		path:   r.path + path,
		router: r.router,
	}
}

func (r *Router) Get(path string, fn HandlerFunc) {
	r.router.GET(r.path+path, toFasthttpHandler(fn))
}

func (r *Router) Post(path string, fn HandlerFunc) {
	r.router.POST(r.path+path, toFasthttpHandler(fn))
}

func (r *Router) Put(path string, fn HandlerFunc) {
	r.router.PUT(r.path+path, toFasthttpHandler(fn))
}

func (r *Router) Patch(path string, fn HandlerFunc) {
	r.router.PATCH(r.path+path, toFasthttpHandler(fn))
}

func (r *Router) Options(path string, fn HandlerFunc) {
	r.router.OPTIONS(r.path+path, toFasthttpHandler(fn))
}

func (r *Router) Head(path string, fn HandlerFunc) {
	r.router.HEAD(r.path+path, toFasthttpHandler(fn))
}

func (r *Router) Delete(path string, fn HandlerFunc) {
	r.router.DELETE(r.path+path, toFasthttpHandler(fn))
}
