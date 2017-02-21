package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/view"
)

type mypage struct {
	Title   string
	Message string
}

func main() {
	app := iris2.New()

	tmpl := view.HTML("./templates", ".html")
	tmpl.Layout("layout.html")

	app.Adapt(tmpl)

	app.Get("/", func(ctx *iris2.Context) {
		ctx.Render("mypage.html", mypage{"My Page title", "Hello world!"}, iris2.Map{"gzip": true})
		// Note that: you can pass "layout" : "otherLayout.html" to bypass the config's Layout property
		// or iris2.NoLayout to disable layout on this render action.
		// third is an optional parameter
	})

	app.Listen(":8080")
}
