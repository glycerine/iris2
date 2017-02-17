package cors

//  +------------------------------------------------------------+
//  | Cors wrapper usage                                         |
//  +------------------------------------------------------------+
//
// import "github.com/go-iris2/iris2/adaptors/cors"
//
// app := iris.New()
// app.Adapt(cors.New(cors.Options{})))

import (
	"github.com/rs/cors"
	"github.com/go-iris2/iris2"
)

// Options is a configuration container to setup the CORS.
type Options cors.Options

// New returns a new cors router wrapper policy
// with the provided options.
//
// create a new cors middleware, pass its options (casted)
// and pass the .ServeHTTP(w,r,next(http.HandlerFunc)) as the wrapper
// for the whole iris' Router
//
// Unlike the cors middleware, this wrapper wraps the entirely router,
// it cannot be registered to specific routes. It works better and all Options are working.
func New(opts Options) iris.RouterWrapperPolicy { return cors.New(cors.Options(opts)).ServeHTTP }

// Default returns a New cors wrapper with the default options
// accepting GET and POST requests and allowing all origins.
func Default() iris.RouterWrapperPolicy { return New(Options{}) }
