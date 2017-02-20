package main

import (
	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/cors"
)

func main() {

	app := iris2.New()
	app.Adapt(iris2.DevLogger())

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	app.Adapt(crs) // this line should be added
	// adaptor supports cors allowed methods, middleware does not.

	// if you want per-route-only cors
	// then you should check https://github.com/iris-contrib/middleware/tree/master/cors

	v1 := app.Party("/api/v1")
	{
		v1.Post("/home", func(c *iris2.Context) {
			c.Log("lalala")
			c.WriteString("Hello from /home")
		})
		v1.Get("/g", func(c *iris2.Context) {
			c.Log("lalala")
			c.WriteString("Hello from /home")
		})
		v1.Post("/h", func(c *iris2.Context) {
			c.Log("lalala")
			c.WriteString("Hello from /home")
		})
	}

	app.Listen(":8080")
}
