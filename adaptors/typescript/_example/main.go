package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/typescript"
)

func main() {
	app := iris2.New()

	ts := typescript.New()
	ts.Config.Dir = "./www/scripts"
	app.Adapt(ts) // adapt the typescript compiler adaptor

	app.StaticWeb("/", "./www") // serve the index.html
	app.Listen(":8080")
}

// open http://localhost:8080
// go to ./www/scripts/app.ts
// make a change
// reload the http://localhost:8080 and you should see the changes
//
// what it does?
// - compiles the typescript files using default compiler options if not tsconfig found
// - watches for changes on typescript files, if a change then it recompiles the .ts to .js
//
// same as you used to do with gulp-like tools, but here at Iris I do my bests to help GO developers.
