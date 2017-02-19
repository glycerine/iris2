package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/middleware/logger"
)

func main() {
	app := iris2.New()

	app.Adapt(iris2.DevLogger()) // it just enables the print of the iris2.DevMode logs. Enable it to view the middleware's messages.

	customLogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
	})

	app.Use(customLogger)

	app.Get("/", func(ctx *iris2.Context) {
		ctx.Writef("hello")
	})

	app.Get("/1", func(ctx *iris2.Context) {
		ctx.Writef("hello")
	})

	app.Get("/2", func(ctx *iris2.Context) {
		ctx.Writef("hello")
	})

	// log http errors
	errorLogger := logger.New()

	app.OnError(iris2.StatusNotFound, func(ctx *iris2.Context) {
		errorLogger.Serve(ctx)
		ctx.Writef("My Custom 404 error page ")
	})

	// http://localhost:8080
	// http://localhost:8080/1
	// http://localhost:8080/2
	app.Listen(":8080")

}
