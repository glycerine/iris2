package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-iris2/iris2"
	"github.com/go-iris2/iris2/adaptors/view"
	"github.com/go-iris2/iris2/adaptors/websocket"
)

type clientPage struct {
	Title string
	Host  string
}

func main() {
	app := iris2.New()
	app.Adapt(view.HTML("./templates", ".html")) // select the html engine to serve templates

	ws := websocket.New(websocket.Config{
		// the path which the websocket client should listen/registered to,
		Endpoint: "/my_endpoint",
		// the client-side javascript static file path
		// which will be served by Iris.
		// default is /iris-ws.js
		// if you change that you have to change the bottom of templates/client.html
		// script tag:
		ClientSourcePath: "/iris-ws.js",
		//
		// Set the timeouts, 0 means no timeout
		// websocket has more configuration, go to ../../config.go for more:
		// WriteTimeout: 0,
		// ReadTimeout:  0,
		// by-default all origins are accepted, you can change this behavior by setting:
		// CheckOrigin: (r *http.Request ) bool {},
		//
		//
		// IDGenerator used to create (and later on, set)
		// an ID for each incoming websocket connections (clients).
		// The request is an argument which you can use to generate the ID (from headers for example).
		// If empty then the ID is generated by DefaultIDGenerator: randomString(64):
		// IDGenerator func(ctx *iris2.Context) string {},
	})

	app.Adapt(ws) // adapt the websocket server, you can adapt more than one with different Endpoint

	app.StaticWeb("/js", "./static/js") // serve our custom javascript code

	app.Get("/", func(ctx *iris2.Context) {
		ctx.Render("client.html", clientPage{"Client Page", ctx.ServerHost()})
	})

	Conn := make(map[websocket.Connection]bool)
	var myChatRoom = "room1"
	var mutex = new(sync.Mutex)

	ws.OnConnection(func(c websocket.Connection) {
		c.Join(myChatRoom)
		mutex.Lock()
		Conn[c] = true
		mutex.Unlock()
		c.On("chat", func(message string) {
			if message == "leave" {
				c.Leave(myChatRoom)
				c.To(myChatRoom).Emit("chat", "Client with ID: "+c.ID()+" left from the room and cannot send or receive message to/from this room.")
				c.Emit("chat", "You have left from the room: "+myChatRoom+" you cannot send or receive any messages from others inside that room.")
				return
			}
		})
		c.OnDisconnect(func() {
			mutex.Lock()
			delete(Conn, c)
			mutex.Unlock()
			fmt.Printf("\nConnection with ID: %s has been disconnected!\n", c.ID())
		})
	})

	var delay = 1 * time.Second
	go func() {
		i := 0
		for {
			mutex.Lock()
			broadcast(Conn, fmt.Sprintf("aaaa %d\n", i))
			mutex.Unlock()
			time.Sleep(delay)
			i++
		}
	}()

	go func() {
		i := 0
		for {
			mutex.Lock()
			broadcast(Conn, fmt.Sprintf("aaaa2 %d\n", i))
			mutex.Unlock()
			time.Sleep(delay)
			i++
		}
	}()

	app.Listen(":8080")
}

func broadcast(Conn map[websocket.Connection]bool, message string) {
	for k := range Conn {
		k.To("room1").Emit("chat", message)
	}
}
