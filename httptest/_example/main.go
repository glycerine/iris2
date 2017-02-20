package main

import (
	"net/http"

	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/sessions"
)

func newApp() *iris2.Framework {
	app := iris2.New()
	app.Adapt(sessions.New(sessions.Config{Cookie: "mysessionid"}))

	app.Get("/hello", func(ctx *iris2.Context) {
		sess := ctx.Session()
		if !sess.HasFlash() /* or sess.GetFlash("name") == "", same thing here */ {
			ctx.HTML(http.StatusUnauthorized, "<h1> Unauthorized Page! </h1>")
			return
		}

		ctx.JSON(http.StatusOK, iris2.Map{
			"Message": "Hello",
			"From":    sess.GetFlash("name"),
		})
	})

	app.Post("/login", func(ctx *iris2.Context) {
		sess := ctx.Session()
		if !sess.HasFlash() {
			sess.SetFlash("name", ctx.FormValue("name"))
		}
		// let's no redirect, just set the flash message, nothing more.
	})

	return app
}

func main() {
	app := newApp()
	app.Listen(":8080")
}
