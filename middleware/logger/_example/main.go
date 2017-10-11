package main

import (
	"net/http"

	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/middleware/logger"
)

func main() {
	app := iris2.New()

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

	app.OnError(http.StatusNotFound, func(ctx *iris2.Context) {
		errorLogger.Serve(ctx)
		ctx.Writef("My Custom 404 error page ")
	})

	// http://localhost:8080
	// http://localhost:8080/1
	// http://localhost:8080/2
	app.Listen(":8080")

}