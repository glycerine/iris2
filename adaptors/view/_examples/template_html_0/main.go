package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/httprouter"
	"github.com/go-iris2/iris2/adaptors/view"
)

func main() {
	app := iris.New(iris.Configuration{Gzip: false, Charset: "UTF-8"}) // defaults to these

	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())

	app.Adapt(view.HTML("./templates", ".html"))

	app.Get("/hi", hi)
	app.Listen(":8080")
}

func hi(ctx *iris.Context) {
	ctx.MustRender("hi.html", struct{ Name string }{Name: "iris"})
}
