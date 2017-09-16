package iris2_test

import (
	"net/http"
	"testing"

	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/httptest"
	"github.com/valyala/fasthttp"
)

func TestBasic(t *testing.T) {
	app := iris2.New()
	app.Get("/", func(ctx *iris2.Context) error {
		ctx.WriteString("Hello World!")
		return nil
	})

	c := httptest.New(app, t)
	c.Get("https://127.0.0.1/").Status(fasthttp.StatusOK).Body("Hello World!")
}

func TestCustomResponseCode(t *testing.T) {
	app := iris2.New()
	app.Get("/", func(ctx *iris2.Context) error {
		ctx.SetStatusCode(123)
		return nil
	})

	c := httptest.New(app, t)
	c.Get("/").Status(123)
}

func TestRoutePath(t *testing.T) {
	app := iris2.New()
	app.Get("/route/path/you/know", func(ctx *iris2.Context) error {
		ctx.WriteString("hi")
		return nil
	})

	c := httptest.New(app, t)
	c.Get("/route/path/you/know").Status(http.StatusOK).Body("hi")
	c.Get("/route/path/you").Status(http.StatusNotFound)
	c.Get("/route/path/you/know/not").Status(http.StatusNotFound)
}

func TestParamPath(t *testing.T) {
	app := iris2.New()
	app.Get("/parstr/:param", func(ctx *iris2.Context) error {
		ctx.WriteString(ctx.ParamString("param"))
		return nil
	})

	c := httptest.New(app, t)
	c.Get("https://127.0.0.1/parstr/%20").Status(http.StatusOK).Body(" ")
	c.Get("https://127.0.0.1/parstr/forreal").Status(http.StatusOK).Body("forreal")
	c.Get("https://127.0.0.1/parstr/for%2Dreal").Status(http.StatusOK).Body("for-real")
}

func TestGroup(t *testing.T) {
	app := iris2.New()
	r := app.Group("/test")
	r.Get("/parstr/:param", func(ctx *iris2.Context) error {
		ctx.WriteString(ctx.ParamString("param"))
		return nil
	})

	c := httptest.New(app, t)
	c.Get("https://127.0.0.1/parstr/%20").Status(http.StatusNotFound)
	c.Get("https://127.0.0.1/test/parstr/%20").Status(http.StatusOK).Body(" ")
}
