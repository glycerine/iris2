package iris2

import (
	"github.com/go-iris2/iris2/session"
	"github.com/valyala/fasthttp"
)

type Context struct {
	*fasthttp.RequestCtx
	rtr  *Router
	sess *session.Session
}

func (c *Context) ParamString(key string) string {
	v := c.UserValue(key)
	s, ok := v.(string)
	if ok {
		return s
	}
	return ""
}

func (c *Context) Session() *session.Session {
	if c.sess == nil {
		c.sess = c.rtr.sessionHandler.Start(c.RequestCtx)
	}
	return c.sess
}

func (c *Context) isSessionStarted() bool {
	return c.sess != nil
}
