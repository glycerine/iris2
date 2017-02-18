package main

import (
	"fmt"

	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/httprouter"
	"github.com/go-iris2/iris2/adaptors/view"
	"github.com/go-iris2/iris2/adaptors/websocket"
)

/* Native messages no need to import the iris-ws.js to the ./templates.client.html
Use of: OnMessage and EmitMessage.


NOTICE: IF YOU HAVE RAN THE PREVIOUS EXAMPLES YOU HAVE TO CLEAR YOUR BROWSER's CACHE
BECAUSE chat.js is different than the CACHED. OTHERWISE YOU WILL GET Ws is undefined from the browser's console, becuase it will use the cached.
*/

type clientPage struct {
	Title string
	Host  string
}

func main() {
	app := iris2.New()
	app.Adapt(iris2.DevLogger())                  // enable all (error) logs
	app.Adapt(httprouter.New())                  // select the httprouter as the servemux
	app.Adapt(view.HTML("./templates", ".html")) // select the html engine to serve templates

	ws := websocket.New(websocket.Config{
		// the path which the websocket client should listen/registered to,
		Endpoint: "/my_endpoint",
		// to enable binary messages (useful for protobuf):
		// BinaryMessages: true,
	})

	app.Adapt(ws) // adapt the websocket server, you can adapt more than one with different Endpoint

	app.StaticWeb("/js", "./static/js") // serve our custom javascript code

	app.Get("/", func(ctx *iris2.Context) {
		ctx.Render("client.html", clientPage{"Client Page", ctx.Host()})
	})

	ws.OnConnection(func(c websocket.Connection) {

		c.OnMessage(func(data []byte) {
			message := string(data)
			c.To(websocket.Broadcast).EmitMessage([]byte("Message from: " + c.ID() + "-> " + message)) // broadcast to all clients except this
			c.EmitMessage([]byte("Me: " + message))                                                    // writes to itself
		})

		c.OnDisconnect(func() {
			fmt.Printf("\nConnection with ID: %s has been disconnected!", c.ID())
		})

	})

	app.Listen(":8080")

}
