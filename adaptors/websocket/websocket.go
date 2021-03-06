// Package websocket provides an easy way to setup server and client side rich websocket experience for Iris
// As originally written by me at https://github.com/kataras/go-websocket
package websocket

import (
	"strings"

	"github.com/go-iris2/iris2"
)

// New returns a new websocket server policy adaptor.
func New(cfg Config) Server {
	return &server{
		config: cfg.Validate(),
		rooms:  make(map[string][]string, 0),
		onConnectionListeners: make([]ConnectionFunc, 0),
	}
}

func fixPath(s string) string {
	if s == "" {
		return ""
	}

	if s[0] != '/' {
		s = "/" + s
	}

	s = strings.Replace(s, "//", "/", -1)
	return s
}

// Adapt implements the iris' adaptor, it adapts the websocket server to an Iris station.
func (s *server) Adapt(frame *iris2.Policies) {
	// bind the server's Handler to Iris at Boot state
	evt := iris2.EventPolicy{
		Boot: func(f *iris2.Framework) {
			wsPath := fixPath(s.config.Endpoint)
			if wsPath == "" {
				f.Log("websocket's configuration field 'Endpoint' cannot be empty, websocket server stops")
				return
			}

			wsClientSidePath := fixPath(s.config.ClientSourcePath)
			if wsClientSidePath == "" {
				f.Log("websocket's configuration field 'ClientSourcePath' cannot be empty, websocket server stops")
				return
			}

			// set the routing for client-side source (javascript) (optional)
			clientSideLookupName := "iris-websocket-client-side"
			wsHandler := s.Handler()
			f.Get(wsPath, wsHandler)
			// check if client side doesn't already exists
			if f.Routes().Lookup(clientSideLookupName) == nil {
				// serve the client side on domain:port/iris-ws.js
				f.StaticContent(wsClientSidePath, "application/javascript", ClientSource).ChangeName(clientSideLookupName)
			}
		},
	}

	evt.Adapt(frame)
}
