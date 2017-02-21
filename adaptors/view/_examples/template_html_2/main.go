package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/view"
)

func main() {
	app := iris2.New()

	tmpl := view.HTML("./templates", ".html")
	tmpl.Layout("layouts/layout.html")
	tmpl.Funcs(map[string]interface{}{
		"greet": func(s string) string {
			return "Greetings " + s + "!"
		},
	})

	app.Adapt(tmpl)

	app.Get("/", func(ctx *iris2.Context) {
		if err := ctx.Render("page1.html", nil); err != nil {
			println(err.Error())
		}
	})

	// remove the layout for a specific route
	app.Get("/nolayout", func(ctx *iris2.Context) {
		if err := ctx.Render("page1.html", nil, iris2.RenderOptions{"layout": iris2.NoLayout}); err != nil {
			println(err.Error())
		}
	})

	// set a layout for a party, .Layout should be BEFORE any Get or other Handle party's method
	my := app.Party("/my").Layout("layouts/mylayout.html")
	{
		my.Get("/", func(ctx *iris2.Context) {
			ctx.MustRender("page1.html", nil)
		})
		my.Get("/other", func(ctx *iris2.Context) {
			ctx.MustRender("page1.html", nil)
		})
	}

	app.Listen(":8080")
}
