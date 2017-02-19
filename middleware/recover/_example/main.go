package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/middleware/recover"
)

func main() {
	app := iris2.New()

	app.Adapt(iris2.DevLogger()) // fast way to enable non-fatal messages to be printed to the user
	// (yes in iris even recover's errors are not fatal because it's restarting,
	// ProdMode messages are only for things that Iris cannot continue at all,
	// these are logged by-default but you can change that behavior too by passing a different LoggerPolicy to the .Adapt)
	app.Use(recover.New()) // it's io.Writer is the same as app.Config.LoggerOut

	i := 0
	// let's simmilate a panic every next request
	app.Get("/", func(ctx *iris2.Context) {
		i++
		if i%2 == 0 {
			panic("a panic here")
		}
		ctx.Writef("Hello, refresh one time more to get panic!")
	})

	// http://localhost:8080
	app.Listen(":8080")
}
