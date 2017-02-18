package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/httprouter"
)

func hello(ctx *iris2.Context) {
	ctx.Writef("Hello from %s", ctx.Path())
}

func main() {
	app := iris2.New()
	app.Adapt(iris2.DevLogger())
	app.Adapt(httprouter.New())

	app.OnError(iris2.StatusNotFound, func(ctx *iris2.Context) {
		ctx.HTML(iris2.StatusNotFound, "<h1>Custom not found handler </h1>")
	})

	app.Get("/", hello)
	app.Get("/users/:userid", func(ctx *iris2.Context) {
		ctx.Writef("Hello user with  id: %s", ctx.Param("userid"))
	})

	app.Get("/myfiles/*file", func(ctx *iris2.Context) {
		ctx.HTML(iris2.StatusOK, "Hello, the dynamic path after /myfiles is:<br/> <b>"+ctx.Param("file")+"</b>")
	})

	app.Get("/users/:userid/messages/:messageid", func(ctx *iris2.Context) {
		ctx.HTML(iris2.StatusOK, `Message from user with id:<br/> <b>`+ctx.Param("userid")+`</b>,
            message id: <b>`+ctx.Param("messageid")+`</b>`)
	})

	// http://127.0.0.1:8080/users/42
	// http://127.0.0.1:8080/myfiles/mydirectory/myfile.zip
	// http://127.0.0.1:8080/users/42/messages/1
	app.Listen(":8080")
}
