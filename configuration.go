package iris2

import (
	"strconv"
	"time"

	"github.com/imdario/mergo"
)

type (
	// OptionSetter sets a configuration field to the main configuration
	// used to help developers to write less and configure only what
	// they really want and nothing else.
	//
	// Usage:
	// iris2.New(iris2.Configuration{Charset: "UTF-8", Gzip:true})
	// now can be done also by using iris2.Option$FIELD:
	// iris2.New(iris2.OptionCharset("UTF-8"), iris2.OptionGzip(true))
	//
	// Benefits:
	// 1. Developers have no worries what option to pass,
	//    they can just type iris2.Option and all options should
	//    be visible to their editor's autocomplete-popup window
	// 2. Can be passed with any order
	// 3. Can override previous configuration
	OptionSetter interface {
		// Set receives a pointer to the global Configuration type and does the job of filling it
		Set(c *Configuration)
	}
	// OptionSet implements the OptionSetter
	OptionSet func(c *Configuration)
)

// Set is the func which makes the OptionSet an OptionSetter, this is used mostly
func (o OptionSet) Set(c *Configuration) {
	o(c)
}

// Configuration the whole configuration for an Iris station instance
// these can be passed via options also, look at the top of this file(configuration.go).
// Configuration is a valid OptionSetter.
type Configuration struct {
	// VHost is the addr or the domain that server listens to, which it's optional
	// When to set VHost manually:
	// 1. it's automatically setted when you're calling
	//     $instance.Listen/ListenUNIX/ListenTLS/ListenLETSENCRYPT functions or
	//     ln,_ := iris2.TCP4/UNIX/TLS/LETSENCRYPT; $instance.Serve(ln)
	// 2. If you using a balancer, or something like nginx
	//    then set it in order to have the correct url
	//    when calling the template helper '{{url }}'
	//    *keep note that you can use {{urlpath }}) instead*
	//
	// Note: this is the main's server Host, you can setup unlimited number of net/http servers
	// listening to the $instance.Handler after the manually-called $instance.Build
	//
	// Default comes from iris2.Default.Listen/.Serve with iris' listeners (iris2.TCP4/UNIX/TLS/LETSENCRYPT).
	VHost string

	// VScheme is the scheme (http:// or https://) putted at the template function '{{url }}'
	// It's an optional field,
	// When to set VScheme manually:
	// 1. You didn't start the main server using $instance.Listen/ListenTLS/ListenLETSENCRYPT
	//    or $instance.Serve($instance.TCP4()/.TLS...)
	// 2. if you're using something like nginx and have iris listening with
	//   addr only(http://) but the nginx mapper is listening to https://
	//
	// Default comes from iris2.Default.Listen/.Serve with iris' listeners (TCP4,UNIX,TLS,LETSENCRYPT).
	VScheme string

	// ReadTimeout is the maximum duration before timing out read of the request.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out write of the response.
	WriteTimeout time.Duration

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int

	// DisablePathCorrection corrects and redirects the requested path to the registered path
	// for example, if /home/ path is requested but no handler for this Route found,
	// then the Router checks if /home handler exists, if yes,
	// (permant)redirects the client to the correct path /home
	//
	// Defaults to false.
	DisablePathCorrection bool

	// EnablePathEscape when is true then its escapes the path, the named parameters (if any).
	// Change to false it if you want something like this https://github.com/kataras/iris/issues/135 to work
	//
	// When do you need to Disable(false) it:
	// accepts parameters with slash '/'
	// Request: http://localhost:8080/details/Project%2FDelta
	// ctx.Param("project") returns the raw named parameter: Project%2FDelta
	// which you can escape it manually with net/url:
	// projectName, _ := url.QueryUnescape(c.Param("project").
	//
	// Defaults to false.
	EnablePathEscape bool

	// FireMethodNotAllowed if it's true router checks for StatusMethodNotAllowed(405) and
	//  fires the 405 error instead of 404
	// Defaults to false.
	FireMethodNotAllowed bool

	// DisableBodyConsumptionOnUnmarshal manages the reading behavior of the context's body readers/binders.
	// If setted to true then it
	// disables the body consumption by the `context.UnmarshalBody/ReadJSON/ReadXML`.
	//
	// By-default io.ReadAll` is used to read the body from the `context.Request.Body which is an `io.ReadCloser`,
	// if this field setted to true then a new buffer will be created to read from and the request body.
	// The body will not be changed and existing data before the
	// context.UnmarshalBody/ReadJSON/ReadXML will be not consumed.
	DisableBodyConsumptionOnUnmarshal bool

	// TimeFormat time format for any kind of datetime parsing
	// Defaults to  "Mon, 02 Jan 2006 15:04:05 GMT".
	TimeFormat string

	// Charset character encoding for various rendering
	// used for templates and the rest of the responses
	// Defaults to "UTF-8".
	Charset string

	// Gzip enables gzip compression on your Render actions, this includes any type of render,
	// templates and pure/raw content
	// If you don't want to enable it globally, you could just use the third parameter
	// on context.Render("myfileOrResponse", structBinding{}, iris2.RenderOptions{"gzip": true})
	// Defaults to false.
	Gzip bool

	// AutoFlashMessage adds the flash-message "msg" automatically to the parameters
	// for each Render, side-effect is that when render is performed, the flash-message
	// shall be garbage-collected!
	AutoFlashMessage bool

	// Other are the custom, dynamic options, can be empty.
	// This field used only by you to set any app's options you want
	// or by custom adaptors, it's a way to simple communicate between your adaptors (if any)
	// Defaults to a non-nil empty map.
	Other map[string]interface{}
}

// Set implements the OptionSetter
func (c Configuration) Set(main *Configuration) {
	if err := mergo.MergeWithOverwrite(main, c); err != nil {
		panic("FATAL ERROR .Configuration as OptionSetter: " + err.Error())
	}
}

// All options starts with "Option" preffix in order to be easier to find what dev searching for
var (

	// OptionVHost is the addr or the domain that server listens to, which it's optional
	// When to set VHost manually:
	// 1. it's automatically setted when you're calling
	//     $instance.Listen/ListenUNIX/ListenTLS/ListenLETSENCRYPT functions or
	//     ln,_ := iris2.TCP4/UNIX/TLS/LETSENCRYPT; $instance.Serve(ln)
	// 2. If you using a balancer, or something like nginx
	//    then set it in order to have the correct url
	//    when calling the template helper '{{url }}'
	//    *keep note that you can use {{urlpath }}) instead*
	//
	// Note: this is the main's server Host, you can setup unlimited number of net/http servers
	// listening to the $instance.Handler after the manually-called $instance.Build
	//
	// Default comes from iris2.Default.Listen/.Serve with iris' listeners (iris2.TCP4/UNIX/TLS/LETSENCRYPT).
	OptionVHost = func(val string) OptionSet {
		return func(c *Configuration) {
			c.VHost = val
		}
	}

	// OptionVScheme is the scheme (http:// or https://) putted at the template function '{{url }}'
	// It's an optional field,
	// When to set Scheme manually:
	// 1. You didn't start the main server using $instance.Listen/ListenTLS/ListenLETSENCRYPT
	//     or $instance.Serve($instance.TCP4()/.TLS...)
	// 2. if you're using something like nginx and have iris listening with
	//    addr only(http://) but the nginx mapper is listening to https://
	//
	// Default comes from iris2.Default.Listen/.Serve with iris' listeners (TCP4,UNIX,TLS,LETSENCRYPT).
	OptionVScheme = func(val string) OptionSet {
		return func(c *Configuration) {
			c.VScheme = val
		}
	}

	// OptionReadTimeout sets the Maximum duration before timing out read of the request.
	OptionReadTimeout = func(val time.Duration) OptionSet {
		return func(c *Configuration) {
			c.ReadTimeout = val
		}
	}

	// OptionWriteTimeout sets the Maximum duration before timing out write of the response.
	OptionWriteTimeout = func(val time.Duration) OptionSet {
		return func(c *Configuration) {
			c.WriteTimeout = val
		}
	}

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes(8MB) is used.
	OptionMaxHeaderBytes = func(val int) OptionSet {
		return func(c *Configuration) {
			c.MaxHeaderBytes = val
		}
	}

	// OptionDisablePathCorrection corrects and redirects the requested path to the registered path
	// for example, if /home/ path is requested but no handler for this Route found,
	// then the Router checks if /home handler exists, if yes,
	// (permant)redirects the client to the correct path /home
	//
	// Defaults to false.
	OptionDisablePathCorrection = func(val bool) OptionSet {
		return func(c *Configuration) {
			c.DisablePathCorrection = val
		}

	}

	// OptionEnablePathEscape when is true then its escapes the path, the named path parameters (if any).
	OptionEnablePathEscape = func(val bool) OptionSet {
		return func(c *Configuration) {
			c.EnablePathEscape = val
		}
	}

	// FireMethodNotAllowed if it's true router checks for StatusMethodNotAllowed(405)
	// and fires the 405 error instead of 404
	// Defaults to false.
	OptionFireMethodNotAllowed = func(val bool) OptionSet {
		return func(c *Configuration) {
			c.FireMethodNotAllowed = val
		}
	}

	// OptionDisableBodyConsumptionOnUnmarshal manages the reading behavior of the context's body readers/binders.
	// If setted to true then it
	// disables the body consumption by the `context.UnmarshalBody/ReadJSON/ReadXML`.
	//
	// By-default io.ReadAll` is used to read the body from the `context.Request.Body which is an `io.ReadCloser`,
	// if this field setted to true then a new buffer will be created to read from and the request body.
	// The body will not be changed and existing data before the context.UnmarshalBody/ReadJSON/ReadXML will be not consumed.
	OptionDisableBodyConsumptionOnUnmarshal = func(val bool) OptionSet {
		return func(c *Configuration) {
			c.DisableBodyConsumptionOnUnmarshal = val
		}
	}

	// OptionTimeFormat time format for any kind of datetime parsing.
	// Defaults to  "Mon, 02 Jan 2006 15:04:05 GMT".
	OptionTimeFormat = func(val string) OptionSet {
		return func(c *Configuration) {
			c.TimeFormat = val
		}
	}

	// OptionCharset character encoding for various rendering
	// used for templates and the rest of the responses
	// Defaults to "UTF-8".
	OptionCharset = func(val string) OptionSet {
		return func(c *Configuration) {
			c.Charset = val
		}
	}

	// OptionGzip enables gzip compression on your Render actions, this includes any type of render, templates and pure/raw content
	// If you don't want to enable it globally, you could just use the third parameter on context.Render("myfileOrResponse", structBinding{}, iris2.RenderOptions{"gzip": true})
	// Defaults to false.
	OptionGzip = func(val bool) OptionSet {
		return func(c *Configuration) {
			c.Gzip = val
		}
	}

	OptionAutoFlashMessage = func(val bool) OptionSet {
		return func(c *Configuration) {
			c.AutoFlashMessage = val
		}
	}

	// Other are the custom, dynamic options, can be empty.
	// This field used only by you to set any app's options you want
	// or by custom adaptors, it's a way to simple communicate between your adaptors (if any)
	// Defaults to a non-nil empty map.
	OptionOther = func(key string, val interface{}) OptionSet {
		return func(c *Configuration) {
			if c.Other == nil {
				c.Other = make(map[string]interface{}, 0)
			}
			c.Other[key] = val
		}
	}
)

var (
	// DefaultTimeFormat default time format for any kind of datetime parsing
	DefaultTimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
	// StaticCacheDuration expiration duration for INACTIVE file handlers, it's a global configuration field to all iris instances
	StaticCacheDuration = 20 * time.Second
)

// Default values for base Iris conf
const (
	DefaultDisablePathCorrection = false
	DefaultEnablePathEscape      = false
	DefaultCharset               = "UTF-8"
	// Per-connection buffer size for requests' reading.
	// This also limits the maximum header size.
	//
	// Increase this buffer if your clients send multi-KB RequestURIs
	// and/or multi-KB headers (for example, BIG cookies).
	//
	// Default buffer size is 8MB
	DefaultMaxHeaderBytes = 8096

	// DefaultReadTimeout no read client timeout
	DefaultReadTimeout = 0
	// DefaultWriteTimeout no serve client timeout
	DefaultWriteTimeout = 0
)

// DefaultConfiguration returns the default configuration for an Iris station, fills the main Configuration
func DefaultConfiguration() Configuration {
	return Configuration{
		VHost:                             "",
		VScheme:                           "",
		ReadTimeout:                       DefaultReadTimeout,
		WriteTimeout:                      DefaultWriteTimeout,
		MaxHeaderBytes:                    DefaultMaxHeaderBytes,
		DisablePathCorrection:             DefaultDisablePathCorrection,
		EnablePathEscape:                  DefaultEnablePathEscape,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		TimeFormat:                        DefaultTimeFormat,
		Charset:                           DefaultCharset,
		Gzip:                              false,
		AutoFlashMessage:                  true,
		Other:                             make(map[string]interface{}, 0),
	}
}

// Default values for base Server conf
const (
	// DefaultServerHostname returns the default hostname which is 0.0.0.0
	DefaultServerHostname = "0.0.0.0"
	// DefaultServerPort returns the default port which is 8080, not used
	DefaultServerPort = 8080
)

var (
	// DefaultServerAddr the default server addr which is: 0.0.0.0:8080
	DefaultServerAddr = DefaultServerHostname + ":" + strconv.Itoa(DefaultServerPort)
)
