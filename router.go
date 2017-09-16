package iris2

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/go-iris2/iris2/session"
	"github.com/valyala/fasthttp"
)

type HandlerFunc func(ctx *Context) error

type Router struct {
	path           string
	sessionHandler *session.Handler
	router         *fasthttprouter.Router
}

func newRouter() *Router {
	return &Router{
		path:           "",
		sessionHandler: session.New(),
		router:         fasthttprouter.New(),
	}
}

func (r *Router) toFasthttpHandler(fn HandlerFunc) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		c := &Context{RequestCtx: ctx, rtr: r}
		fn(c)
		if c.isSessionStarted() {
			c.Session().Save()
		}
	}
	// TODO error handling
}

func (r *Router) Group(path string) *Router {
	return &Router{
		path:           r.path + path,
		sessionHandler: r.sessionHandler,
		router:         r.router,
	}
}

func (r *Router) Get(path string, fn HandlerFunc) {
	r.router.GET(r.path+path, r.toFasthttpHandler(fn))
}

func (r *Router) Post(path string, fn HandlerFunc) {
	r.router.POST(r.path+path, r.toFasthttpHandler(fn))
}

func (r *Router) Put(path string, fn HandlerFunc) {
	r.router.PUT(r.path+path, r.toFasthttpHandler(fn))
}

func (r *Router) Patch(path string, fn HandlerFunc) {
	r.router.PATCH(r.path+path, r.toFasthttpHandler(fn))
}

func (r *Router) Options(path string, fn HandlerFunc) {
	r.router.OPTIONS(r.path+path, r.toFasthttpHandler(fn))
}

func (r *Router) Head(path string, fn HandlerFunc) {
	r.router.HEAD(r.path+path, r.toFasthttpHandler(fn))
}

func (r *Router) Delete(path string, fn HandlerFunc) {
	r.router.DELETE(r.path+path, r.toFasthttpHandler(fn))
}
