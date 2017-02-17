package sessions

import (
	"encoding/base64"
	"time"
)

const (
	// DefaultCookieName the secret cookie's name for sessions
	DefaultCookieName = "gosessionsid"
	// DefaultCookieLength is the default Session Manager's CookieLength, which is 32
	DefaultCookieLength = 32
)

// Config the configuration for sessions
// has 5 fields
// first is the cookieName, the session's name (string) ["mysessionsecretcookieid"]
// second enable if you want to decode the cookie's key also
// third is the time which the client's cookie expires
// forth is the cookie length (sessionid) int, defaults to 32, do not change if you don't have any reason to do
// fifth is the DisableSubdomainPersistence which you can set it to true in order dissallow your q subdomains to have access to the session cook
type (
	// OptionSetter used as the type of return of a func which sets a configuration field's value
	OptionSetter interface {
		// Set receives a pointer to the Config type and does the job of filling it
		Set(c *Config)
	}
	// OptionSet implements the OptionSetter
	OptionSet func(c *Config)

	Config struct {
		// Cookie string, the session's client cookie name, for example: "mysessionid"
		//
		// Defaults to "gosessionid"
		Cookie string

		// DecodeCookie set it to true to decode the cookie key with base64 URLEncoding
		//
		// Defaults to false
		DecodeCookie bool

		// Expires the duration of which the cookie must expires (created_time.Add(Expires)).
		// If you want to delete the cookie when the browser closes, set it to -1 but in this case, the server side's session duration is up to GcDuration
		//
		// 0 means no expire, (24 years)
		// -1 means when browser closes
		// > 0 is the time.Duration which the session cookies should expire.
		//
		// Defaults to infinitive/unlimited life duration(0)
		Expires time.Duration

		// CookieLength the length of the sessionid's cookie's value, let it to 0 if you don't want to change it
		//
		// Defaults to 32
		CookieLength int

		// DisableSubdomainPersistence set it to true in order dissallow your q subdomains to have access to the session cookie
		//
		// Defaults to false
		DisableSubdomainPersistence bool
	}
)

// Set implements the OptionSetter
func (c Config) Set(main *Config) {
	c = c.Validate()
	*main = c
}

// Set is the func which makes the OptionSet an OptionSetter, this is used mostly
func (o OptionSet) Set(c *Config) {
	o(c)
}

// Cookie string, the session's client cookie name, for example: "mysessionid"
//
// Defaults to "gosessionid"
func Cookie(val string) OptionSet {
	return func(c *Config) {
		c.Cookie = val
	}
}

// DecodeCookie set it to true to decode the cookie key with base64 URLEncoding
//
// Defaults to false
func DecodeCookie(val bool) OptionSet {
	return func(c *Config) {
		c.DecodeCookie = val
	}
}

// Expires the duration of which the cookie must expires (created_time.Add(Expires)).
// If you want to delete the cookie when the browser closes, set it to -1 but in this case, the server side's session duration is up to GcDuration
//
// Defaults to infinitive/unlimited life duration(0)
func Expires(val time.Duration) OptionSet {
	return func(c *Config) {
		c.Expires = val
	}
}

// CookieLength the length of the sessionid's cookie's value, let it to 0 if you don't want to change it
//
// Defaults to 32
func CookieLength(val int) OptionSet {
	return func(c *Config) {
		c.CookieLength = val
	}
}

// DisableSubdomainPersistence set it to true in order dissallow your q subdomains to have access to the session cookie
//
// Defaults to false
func DisableSubdomainPersistence(val bool) OptionSet {
	return func(c *Config) {
		c.DisableSubdomainPersistence = val
	}
}

// Validate corrects missing fields configuration fields and returns the right configuration
func (c Config) Validate() Config {

	if c.Cookie == "" {
		c.Cookie = DefaultCookieName
	}

	if c.DecodeCookie {
		c.Cookie = base64.URLEncoding.EncodeToString([]byte(c.Cookie)) // change the cookie's name/key to a more safe(?)
		// get the real value for your tests by:
		//sessIdKey := url.QueryEscape(base64.URLEncoding.EncodeToString([]byte(Sessions.Cookie)))
	}

	if c.CookieLength <= 0 {
		c.CookieLength = DefaultCookieLength
	}

	return c
}
