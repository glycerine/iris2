package main

import (
	"time"

	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/middleware/basicauth"
)

func main() {
	app := iris2.New()

	authConfig := basicauth.Config{
		Users:      map[string]string{"myusername": "mypassword", "mySecondusername": "mySecondpassword"},
		Realm:      "Authorization Required", // defaults to "Authorization Required"
		ContextKey: "mycustomkey",            // defaults to "user"
		Expires:    time.Duration(30) * time.Minute,
	}

	authentication := basicauth.New(authConfig)
	app.Get("/", func(ctx *iris2.Context) { ctx.Redirect("/admin") })
	// to global app.Use(authentication) (or app.UseGlobal before the .Listen)
	// to routes
	/*
		app.Get("/mysecret", authentication, func(ctx *iris2.Context) {
			username := ctx.GetString("mycustomkey") //  the Contextkey from the authConfig
			ctx.Writef("Hello authenticated user: %s ", username)
		})
	*/

	// to party

	needAuth := app.Party("/admin", authentication)
	{
		//http://localhost:8080/admin
		needAuth.Get("/", func(ctx *iris2.Context) {
			username := ctx.GetString("mycustomkey") //  the Contextkey from the authConfig
			ctx.Writef("Hello authenticated user: %s from: %s ", username, ctx.Path())
		})
		// http://localhost:8080/admin/profile
		needAuth.Get("/profile", func(ctx *iris2.Context) {
			username := ctx.GetString("mycustomkey") //  the Contextkey from the authConfig
			ctx.Writef("Hello authenticated user: %s from: %s ", username, ctx.Path())
		})
		// http://localhost:8080/admin/settings
		needAuth.Get("/settings", func(ctx *iris2.Context) {
			username := authConfig.User(ctx) // shortcut for ctx.GetString("mycustomkey")
			ctx.Writef("Hello authenticated user: %s from: %s ", username, ctx.Path())
		})
	}

	// open http://localhost:8080/admin
	app.Listen(":8080")
}
