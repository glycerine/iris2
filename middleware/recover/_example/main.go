package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/middleware/recover"
)

func main() {
	app := iris2.New()

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
