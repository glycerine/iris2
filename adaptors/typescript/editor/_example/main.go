package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/httprouter"
	"github.com/go-iris2/iris2/adaptors/typescript" // optinally
	"github.com/go-iris2/iris2/adaptors/typescript/editor"
)

func main() {
	app := iris2.New()
	app.Adapt(iris2.DevLogger())
	app.Adapt(httprouter.New()) // adapt a router, order doesn't matters

	// optionally but good to have, I didn't put inside editor or the editor in the typescript compiler adaptors
	// because you may use tools like gulp and you may use the editor without the typescript compiler adaptor.
	// but if you need auto-compilation on .ts, we have a solution:
	ts := typescript.New()
	ts.Config.Dir = "./www/scripts/"
	app.Adapt(ts) // adapt the typescript compiler adaptor

	editorConfig := editor.Config{
		Hostname:   "127.0.0.1",
		Port:       4444,
		WorkingDir: "./www/scripts/", // "/path/to/the/client/side/directory/",
		Username:   "myusername",
		Password:   "mypassword",
	}
	e := editor.New(editorConfig)
	app.Adapt(e) // adapt the editor

	app.StaticWeb("/", "./www") // serve the index.html

	app.Listen(":8080")
}
