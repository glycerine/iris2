package main

import (
	"time"

	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/sessions"
)

func main() {
	app := iris2.New()

	mySessions := sessions.New(sessions.Config{
		// Cookie string, the session's client cookie name, for example: "mysessionid"
		//
		// Defaults to "irissessionid"
		Cookie: "mysessionid",
		// it's time.Duration, from the time cookie is created, how long it can be alive?
		// 0 means no expire.
		// -1 means expire when browser closes
		// or set a value, like 2 hours:
		Expires: time.Hour * 2,
		// the length of the sessionid's cookie's value
		CookieLength: 32,
		// if you want to invalid cookies on different subdomains
		// of the same host, then enable it
		DisableSubdomainPersistence: false,
	})

	// OPTIONALLY:
	// import "github.com/go-iris2/iris2/adaptors/sessions/sessiondb/redis"
	// or import "github.com/go-iris2/iris2/sessions/sessiondb/$any_available_community_database"
	// mySessions.UseDatabase(redis.New(...))

	app.Adapt(mySessions) // Adapt the session manager we just created.

	app.Get("/", func(ctx *iris2.Context) {
		ctx.Writef("You should navigate to the /set, /get, /delete, /clear,/destroy instead")
	})
	app.Get("/set", func(ctx *iris2.Context) {

		//set session values
		ctx.Session().Set("name", "iris")

		//test if setted here
		ctx.Writef("All ok session setted to: %s", ctx.Session().GetString("name"))
	})

	app.Get("/get", func(ctx *iris2.Context) {
		// get a specific key, as string, if no found returns just an empty string
		name := ctx.Session().GetString("name")

		ctx.Writef("The name on the /set was: %s", name)
	})

	app.Get("/delete", func(ctx *iris2.Context) {
		// delete a specific key
		ctx.Session().Delete("name")
	})

	app.Get("/clear", func(ctx *iris2.Context) {
		// removes all entries
		ctx.Session().Clear()
	})

	app.Get("/destroy", func(ctx *iris2.Context) {

		//destroy, removes the entire session and cookie
		ctx.SessionDestroy()
		msg := "You have to refresh the page to completely remove the session (browsers works this way, it's not iris-specific.)"

		ctx.Writef(msg)
		ctx.Log(msg)
	}) // Note about destroy:
	//
	// You can destroy a session outside of a handler too, using the:
	// mySessions.DestroyByID
	// mySessions.DestroyAll

	app.Listen(":8080")
}
