package iris2

import (
	"github.com/valyala/fasthttp"
)

type Server struct {
	*fasthttp.Server
	*Router
}

func New() *Server {
	s := &Server{
		Server: &fasthttp.Server{},
		Router: newRouter(),
	}
	s.Handler = s.Router.router.Handler
	return s
}
