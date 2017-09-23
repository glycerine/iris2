package iris2

import (
	"fmt"
	"strconv"

	"github.com/go-iris2/iris2/session"
	"github.com/valyala/fasthttp"
)

type Context struct {
	*fasthttp.RequestCtx
	rtr  *Router
	sess *session.Session
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

func (c *Context) Redirect(r string) {
	c.RequestCtx.Redirect(r, fasthttp.StatusTemporaryRedirect)
}

func (c *Context) Method() string {
	return string(c.RequestCtx.Method())
}

func (c *Context) ParamString(key string) string {
	v := c.UserValue(key)
	s, ok := v.(string)
	if ok {
		return s
	}
	return ""
}

func (c *Context) ParamInt(key string) (int, error) {
	v := c.UserValue(key)
	switch i := v.(type) {
	case int:
	case uint:
	case int64:
	case int32:
	case int16:
	case int8:
	case uint16:
	case uint32:
	case uint64:
	case uint8:
	case float32:
	case float64:
		return int(i), nil
	case string:
	case []byte:
		return strconv.Atoi(string(i))
	}
	return -1, fmt.Errorf("key %s is not an number", key)
}
