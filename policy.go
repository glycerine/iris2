package iris2

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-iris2/iris2/errors"
)

type (
	// Policy is an interface which should be implemented by all
	// modules that can adapt a policy to the Framework.
	// With a Policy you can change the behavior of almost each of the existing Iris' features.
	Policy interface {
		// Adapt receives the main *Policies which the Policy should be attached on.
		Adapt(frame *Policies)
	}

	// Policies is the main policies list, the rest of the objects that implement the Policy
	// are adapted to the object which contains a field of type *Policies.
	//
	// Policies can have nested policies behaviors too.
	// See iris2.go field: 'policies' and function 'Adapt' for more.
	Policies struct {
		EventPolicy
		RouterReversionPolicy
		RouterBuilderPolicy
		RouterWrapperPolicy
		RenderPolicy
		TemplateFuncsPolicy
		SessionsPolicy
	}
)

// Adapt implements the behavior in order to be valid to pass Policies as one
// useful for third-party libraries which can provide more tools in one registration.
func (p Policies) Adapt(frame *Policies) {

	// Adapt the flow callbacks (optionally)
	p.EventPolicy.Adapt(frame)

	// Adapt the reverse routing behaviors and policy
	p.RouterReversionPolicy.Adapt(frame)

	// Adapt the router builder
	if p.RouterBuilderPolicy != nil {
		p.RouterBuilderPolicy.Adapt(frame)
	}

	// Adapt any Router's wrapper (optionally)
	if p.RouterWrapperPolicy != nil {
		p.RouterWrapperPolicy.Adapt(frame)
	}

	// Adapt the render policy (both templates and rich content)
	if p.RenderPolicy != nil {
		p.RenderPolicy.Adapt(frame)
	}

	// Adapt the template funcs which can be used to register template funcs
	// from community's packages, it doesn't matters what template/view engine the user
	// uses, and if uses at all.
	if p.TemplateFuncsPolicy != nil {
		p.TemplateFuncsPolicy.Adapt(frame)
	}

	p.SessionsPolicy.Adapt(frame)

}

type (
	// EventListener is the signature for type of func(*Framework),
	// which is used to register events inside an EventPolicy.
	//
	// Keep note that, inside the policy this is a wrapper
	// in order to register more than one listener without the need of slice.
	EventListener func(*Framework)

	// EventPolicy contains the available Framework's flow event callbacks.
	// Available events:
	// - Boot
	// - Build
	// - Interrupted
	// - Recover
	EventPolicy struct {
		// Boot with a listener type of EventListener.
		//   Fires when '.Boot' is called (by .Serve functions or manually),
		//   before the Build of the components and the Listen,
		//   after VHost and VSCheme configuration has been setted.
		Boot EventListener
		// Before Listen, after Boot
		Build EventListener
		// Interrupted with a listener type of EventListener.
		//   Fires after the terminal is interrupted manually by Ctrl/Cmd + C
		//   which should be used to release external resources.
		// Iris will close and os.Exit at the end of custom interrupted events.
		// If you want to prevent the default behavior just block on the custom Interrupted event.
		Interrupted EventListener
		// Recover with a listener type of func(*Framework, interface{}).
		//   Fires when an unexpected error(panic) is happening at runtime,
		//   while the server's net.Listener accepting requests
		//   or when a '.Must' call contains a filled error.
		//   Used to release external resources and '.Close' the server.
		//   Only one type of this callback is allowed.
		//
		//   If not empty then the Framework will skip its internal
		//   server's '.Close' and panic to its '.Logger' and execute that callback instaed.
		//   Differences from Interrupted:
		//    1. Fires on unexpected errors
		//    2. Only one listener is allowed.
		Recover func(*Framework, error)
	}
)

var _ Policy = EventPolicy{}

// Adapt adaps an EventPolicy object to the main *Policies.
func (e EventPolicy) Adapt(frame *Policies) {

	// Boot event listener, before the build (old: PreBuild)
	frame.EventPolicy.Boot =
		wrapEvtListeners(frame.EventPolicy.Boot, e.Boot)

		// Build event listener, after Boot and before Listen(old: PostBuild & PreListen)
	frame.EventPolicy.Build =
		wrapEvtListeners(frame.EventPolicy.Build, e.Build)

		// Interrupted event listener, when control+C or manually interrupt by os signal
	frame.EventPolicy.Interrupted =
		wrapEvtListeners(frame.EventPolicy.Interrupted, e.Interrupted)

	// Recover event listener, when panic on .Must and inside .Listen/ListenTLS/ListenUNIX/ListenLETSENCRYPT/Serve
	// only one allowed, no wrapper is used.
	if e.Recover != nil {
		frame.EventPolicy.Recover = e.Recover
	}

}

// Fire fires an EventListener with its Framework when listener is not nil.
// Returns true when fired, otherwise false.
func (e EventPolicy) Fire(ln EventListener, s *Framework) bool {
	if ln != nil {
		ln(s)
		return true
	}
	return false
}

func wrapEvtListeners(prev EventListener, next EventListener) EventListener {
	if next == nil {
		return prev
	}
	listener := next
	if prev != nil {
		listener = func(s *Framework) {
			prev(s)
			next(s)
		}
	}

	return listener
}

type (
	// RouterReversionPolicy is used for the reverse routing feature on
	// which custom routers should create and adapt to the Policies.
	RouterReversionPolicy struct {
		// StaticPath should return the static part of the route path
		// for example, with the default router (: and *):
		// /api/user/:userid should return /api/user
		// /api/user/:userid/messages/:messageid should return /api/user
		// /dynamicpath/*path should return /dynamicpath
		// /my/path should return /my/path
		StaticPath func(path string) string
		// WildcardPath should return a path converted to a 'dynamic' path
		// for example, with the default router(wildcard symbol: '*'):
		// ("/static", "path") should return /static/*path
		// ("/myfiles/assets", "anything") should return /myfiles/assets/*anything
		WildcardPath func(path string, paramName string) string
		// URLPath used for reverse routing on templates with {{ url }} and {{ path }} funcs.
		// Receives the route name and  arguments and returns its http path
		URLPath func(r RouteInfo, args ...string) string
	}
	// RouterBuilderPolicy is the most useful Policy for custom routers.
	// A custom router should adapt this policy which is a func
	// accepting a route repository (contains all necessary routes information)
	// and a context pool which should be used inside router's handlers.
	RouterBuilderPolicy func(repo RouteRepository, cPool ContextPool) http.Handler
	// RouterWrapperPolicy is the Policy which enables a wrapper on the top of
	// the builded Router. Usually it's useful for third-party middleware
	// when need to wrap the entire application with a middleware like CORS.
	RouterWrapperPolicy func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
)

func normalizePath(path string) string {
	path = strings.Replace(path, "//", "/", -1)
	if len(path) > 1 && strings.IndexByte(path, '/') == len(path)-1 {
		// if  it's not "/" and ending with slash remove that slash
		path = path[0 : len(path)-2]
	}
	return path
}

// Adapt adaps a RouterReversionPolicy object to the main *Policies.
func (r RouterReversionPolicy) Adapt(frame *Policies) {
	if r.StaticPath != nil {
		staticPathFn := r.StaticPath
		frame.RouterReversionPolicy.StaticPath = func(path string) string {
			return staticPathFn(normalizePath(path))
		}
	}

	if r.WildcardPath != nil {
		wildcardPathFn := r.WildcardPath
		frame.RouterReversionPolicy.WildcardPath = func(path string, paramName string) string {
			return wildcardPathFn(normalizePath(path), paramName)
		}
	}

	if r.URLPath != nil {
		frame.RouterReversionPolicy.URLPath = r.URLPath
	}
}

// Adapt adaps a RouterBuilderPolicy object to the main *Policies.
func (r RouterBuilderPolicy) Adapt(frame *Policies) {
	// What is this kataras?
	// The whole design of this file is brilliant = go's power + my ideas and experience on software architecture.
	//
	// When the router decides to compile/build this behavior
	// then this overload will check for a wrapper too
	// if a wrapper exists it will wrap the result of the RouterBuilder (which is http.Handler, the Router.)
	// and return that instead.
	// I moved the logic here so we don't need a 'compile/build' method inside the routerAdaptor.
	frame.RouterBuilderPolicy = RouterBuilderPolicy(func(repo RouteRepository, cPool ContextPool) http.Handler {
		handler := r(repo, cPool)
		wrapper := frame.RouterWrapperPolicy
		if wrapper != nil {
			originalHandler := handler.ServeHTTP

			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wrapper(w, r, originalHandler)
			})
		}
		return handler
	})
}

// Adapt adaps a RouterWrapperPolicy object to the main *Policies.
func (r RouterWrapperPolicy) Adapt(frame *Policies) {
	frame.RouterWrapperPolicy = r
}

// RenderPolicy is the type which you can adapt custom renderers
// based on the 'name', simple as that.
// Note that the whole template view system and
// content negotiation works by setting this function via other adaptors.
//
// The functions are wrapped, like any other policy func, the only difference is that
// here the developer has a priority over the defaults:
//  - the last registered is trying to be executed first
//  - the first registered is executing last.
// So a custom adaptor that the community can create and share with each other
// can override the existing one with just a simple registration.
type RenderPolicy func(out io.Writer, name string, bind interface{}, options ...map[string]interface{}) (error, bool)

// Adapt adaps a RenderPolicy object to the main *Policies.
func (r RenderPolicy) Adapt(frame *Policies) {
	if r != nil {
		renderer := r
		prevRenderer := frame.RenderPolicy
		if prevRenderer != nil {
			nextRenderer := r
			renderer = func(out io.Writer, name string, binding interface{}, options ...map[string]interface{}) (error, bool) {
				// Remember: RenderPolicy works in the opossite order of declaration,
				// the last registered is trying to be executed first,
				// the first registered is executing last.
				err, ok := nextRenderer(out, name, binding, options...)
				if !ok {

					prevErr, prevOk := prevRenderer(out, name, binding, options...)
					if err != nil {
						if prevErr != nil {
							err = errors.New(prevErr.Error()).Append(err.Error())
						}
					}
					if prevOk {
						ok = true
					}
				}
				// this renderer is responsible for this name
				// but it has an error, so don't continue to the next
				return err, ok

			}
		}

		frame.RenderPolicy = renderer
	}
}

// TemplateFuncsPolicy sets or overrides template func map.
// Defaults are the iris2.URL and iris2.Path, all the template engines supports the following:
// {{ url "mynamedroute" "pathParameter_ifneeded"} }
// {{ urlpath "mynamedroute" "pathParameter_ifneeded" }}
// {{ render "header.html" }}
// {{ render_r "header.html" }} // partial relative path to current page
// {{ yield }}
// {{ current }}
//
// Developers can already set the template's func map from the view adaptors, example: view.HTML(...).Funcs(...)),
// this type exists in order to be possible from third-party developers to create packages that bind template functions
// to the Iris without the need of knowing what template engine is used by the user or
// what order of declaration the user should follow.
type TemplateFuncsPolicy map[string]interface{} // interface can be: func(arguments ...string) string {}

// Adapt adaps a TemplateFuncsPolicy object to the main *Policies.
func (t TemplateFuncsPolicy) Adapt(frame *Policies) {
	if len(t) > 0 {
		if frame.TemplateFuncsPolicy == nil {
			frame.TemplateFuncsPolicy = t
			return
		}

		if frame.TemplateFuncsPolicy != nil {
			for k, v := range t {
				// set or replace the existing
				frame.TemplateFuncsPolicy[k] = v
			}
		}
	}
}

type (
	// Author's notes:
	// session manager can work as a middleware too
	// but we want an easy-api for the user
	// as we did before with: context.Session().Set/Get...
	// these things cannot be done with middleware and sessions is a critical part of an application
	// which needs attention, so far we used the kataras/go-sessions which I spent many weeks to create
	// and that time has not any known bugs or any other issues, it's fully featured.
	// BUT user may want to use other session library and in the same time users should be able to use
	// iris' api for sessions from context, so a policy is that we need, the policy will contains
	// the Start(responsewriter, request) and the Destroy(responsewriter, request)
	// (keep note that this Destroy is not called at the end of a handler, Start does its job without need to end something
	// sessions are setting in real time, when the user calls .Set ),
	// the Start(responsewriter, request) will return a 'Session' which will contain the API for context.Session() , it should be
	// rich, as before, so the interface will be a clone of the kataras/go-sessions/Session.
	// If the user wants to use other library and that library missing features that kataras/go-sesisons has
	// then the user should make an empty implementation of these calls in order to work.
	// That's no problem, before they couldn't adapt any session manager, now they will can.
	//
	// The databases or stores registration will be in the session manager's responsibility,
	// as well the DestroyByID and DestroyAll (I'm calling these with these names because
	//                                        I take as base the kataras/go-sessions,
	//                                        I have no idea if other session managers
	//                                        supports these things, if not then no problem,
	//                                        these funcs will be not required by the sessions policy)
	//
	// ok let's begin.

	// Session should expose the SessionsPolicy's end-user API.
	// This will be returned at the sess := context.Session().
	Session interface {
		ID() string
		Get(string) interface{}
		GetString(key string) string
		GetStructure(key string, value interface{}) error
		GetInt(key string) (int, error)
		GetUint(key string) (uint, error)
		GetInt64(key string) (int64, error)
		GetFloat32(key string) (float32, error)
		GetFloat64(key string) (float64, error)
		GetBoolean(key string) (bool, error)
		GetAll() map[string]interface{}
		VisitAll(cb func(k string, v interface{}))
		Set(string, interface{})

		HasFlash() bool
		GetFlash(string) interface{}
		GetFlashString(string) string
		GetFlashes() map[string]interface{}
		SetFlash(string, interface{})

		Delete(string)
		DeleteFlash(string)
		Clear()
		ClearFlashes()
	}

	// SessionsPolicy is the policy for a session manager.
	//
	// A SessionsPolicy should be responsible to Start a sesion based
	// on raw http.ResponseWriter and http.Request, which should return
	// a compatible iris2.Session interface, type. If the external session manager
	// doesn't qualifies, then the user should code the rest of the functions with empty implementation.
	//
	// A SessionsPolicy should be responsible to Destroy a session based
	// on the http.ResponseWriter and http.Request, this function should works individually.
	//
	// No iris2.Context required from users. In order to be able to adapt any external session manager.
	//
	// The SessionsPolicy should be adapted once.
	SessionsPolicy struct {
		// Start should starts the session for the particular net/http request
		Start func(http.ResponseWriter, *http.Request) Session

		// Destroy should kills the net/http session and remove the associated cookie
		// Keep note that: Destroy should not called at the end of any handler, it's an independent func.
		// Start should set
		// the values at realtime and if manager doesn't supports these
		// then the user manually have to call its 'done' func inside the handler.
		Destroy func(http.ResponseWriter, *http.Request)
	}
)

// Adapt adaps a SessionsPolicy object to the main *Policies.
//
// Remember: Each policy is an adaptor.
// An adaptor should contains one or more policies too.
func (s SessionsPolicy) Adapt(frame *Policies) {
	if s.Start != nil {
		frame.SessionsPolicy.Start = s.Start
	}
	if s.Destroy != nil {
		frame.SessionsPolicy.Destroy = s.Destroy
	}
}
