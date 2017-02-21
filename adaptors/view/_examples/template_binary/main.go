package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/view"
)

func main() {
	app := iris2.New()

	//$ go-bindata ./templates/...
	// templates are not used, you can delete the folder and run the example
	app.Adapt(view.HTML("./templates", ".html").Binary(Asset, AssetNames))

	app.Get("/hi", hi)
	app.Listen(":8080")
}

func hi(ctx *iris2.Context) {
	ctx.MustRender("hi.html", struct{ Name string }{Name: "iris"})
}
