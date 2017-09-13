package iris2

import (
	"github.com/valyala/fasthttp"
)

type Context struct {
	*fasthttp.RequestCtx
}

func toFasthttpHandler(fn HandlerFunc) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fn(&Context{RequestCtx: ctx})
	}
	// TODO error handling
}
