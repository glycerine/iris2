package main

import (
	"net/http"

	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/template/amber"
	"github.com/go-iris2/iris2/template/html"
)

type mypage struct {
	Title   string
	Message string
}

// Iris examples covers the most part, including all 6 template engines and their configurations:
// https://github.com/iris-contrib/examples/tree/master/template_engines

func main() {

	iris2.UseTemplate(html.New()) // the Iris' default if no template engines are setted.

	// add our second template engine with the same directory but with .amber file extension
	iris2.UseTemplate(amber.New(amber.Config{})).Directory("./templates", ".amber")

	iris2.Get("/render_html", func(ctx *iris2.Context) {
		ctx.RenderWithStatus(http.StatusOK, "hiHTML.html", map[string]interface{}{"Name": "You!"})
	})

	iris2.Get("/render_amber", func(ctx *iris2.Context) {
		ctx.MustRender("hiAMBER.amber", map[string]interface{}{"Name": "You!"})
	})

	println("Open a browser tab & go to localhost:8080/render_html  & localhost:8080/render_amber")
	iris2.Listen(":8080")
}

// Iris examples covers the most part, including all 6 template engines and their configurations:
// https://github.com/iris-contrib/examples/tree/master/template_engines
