package main

import (
	"encoding/xml"

	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/view"
)

// ExampleXML just a test struct to view represents xml content-type
type ExampleXML struct {
	XMLName xml.Name `xml:"example"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

func main() {
	app := iris2.New()

	app.Get("/data", func(ctx *iris2.Context) {
		ctx.Data(http.StatusOK, []byte("Some binary data here."))
	})

	app.Get("/text", func(ctx *iris2.Context) {
		ctx.Text(http.StatusOK, "Plain text here")
	})

	app.Get("/json", func(ctx *iris2.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"hello": "json"}) // or myjsonStruct{hello:"json}
	})

	app.Get("/jsonp", func(ctx *iris2.Context) {
		ctx.JSONP(http.StatusOK, "callbackName", map[string]string{"hello": "jsonp"})
	})

	app.Get("/xml", func(ctx *iris2.Context) {
		ctx.XML(http.StatusOK, ExampleXML{One: "hello", Two: "xml"}) // or iris2.Map{"One":"hello"...}
	})

	app.Get("/markdown", func(ctx *iris2.Context) {
		ctx.Markdown(http.StatusOK, "# Hello Dynamic Markdown Iris")
	})

	app.Adapt(view.HTML("./templates", ".html"))
	app.Get("/template", func(ctx *iris2.Context) {

		ctx.MustRender(
			"hi.html",                 // the file name of the template relative to the './templates'
			iris2.Map{"Name": "Iris"}, // the .Name inside the ./templates/hi.html
			iris2.Map{"gzip": false},  // enable gzip for big files
		)

	})

	// ------ first customization without even the need of *Context or a Handler--------
	//
	// Custom new content-/type:
	//	app.Adapt(iris2.RenderPolicy(func(out io.Writer, name string, binding interface{}, options ...map[string]interface{}) (error, bool) {
	//		if name == "customcontent-type" {
	//
	//			// some very advanced things here:
	//			out.Write([]byte(binding.(string)))
	//			return nil, true
	//		}
	//		return nil, false
	//	}))
	//
	//	app.Get("/custom", func(ctx *iris2.Context) {
	//		ctx.RenderWithStatus(http.StatusOK, // or MustRender
	//			"customcontent-type",
	//			"my custom content here!",
	//		)
	//	})
	//
	// ---- second -----------------------------------------------------------------------
	//
	// Override the defaults (the json,xml,jsonp,text,data and so on), an existing  content-type:
	//	app.Adapt(iris2.RenderPolicy(func(out io.Writer, name string, binding interface{}, options ...map[string]interface{}) (error, bool) {
	//		if name == "text/plain" {
	//			out.Write([]byte("From the custom text/plain renderer: " + binding.(string)))
	//			return nil, true
	//		}
	//
	//		return nil, false
	//	}))
	// // the context.Text's behaviors was changed now by your custom renderer.
	//

	app.Listen(":8080")
}
