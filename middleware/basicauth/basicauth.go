package basicauth

import (
	"encoding/base64"
	"strconv"
	"time"

	"github.com/go-iris2/iris2"
)

//  +------------------------------------------------------------+
//  | Middleware usage                                           |
//  +------------------------------------------------------------+
//
// import "github.com/go-iris2/iris2/middleware/basicauth"
//
// app := iris2.New()
// authentication := basicauth.Default(map[string]string{"myusername": "mypassword", "mySecondusername": "mySecondpassword"})
// app.Get("/dashboard", authentication, func(ctx *iris2.Context){})
//
// for more configuration basicauth.New(basicauth.Config{...})
// see _example

type (
	encodedUser struct {
		HeaderValue string
		Username    string
		logged      bool
		expires     time.Time
	}
	encodedUsers []encodedUser

	basicAuthMiddleware struct {
		config Config
		// these are filled from the config.Users map at the startup
		auth             encodedUsers
		realmHeaderValue string
		expireEnabled    bool // if the config.Expires is a valid date, default disabled
	}
)

//

// New takes one parameter, the Config returns a HandlerFunc
// use: iris2.UseFunc(New(...)), iris2.Get(...,New(...),...)
func New(c Config) iris2.HandlerFunc {
	b := &basicAuthMiddleware{config: DefaultConfig().MergeSingle(c)}
	b.init()
	return b.Serve
}

// Default takes one parameter, the users returns a HandlerFunc
// use: iris2.UseFunc(Default(...)), iris2.Get(...,Default(...),...)
func Default(users map[string]string) iris2.HandlerFunc {
	c := DefaultConfig()
	c.Users = users
	return New(c)
}

//

// User returns the user from context key same as 'ctx.GetString("user")' but cannot be used by the developer, use the basicauth.Config.User func instead.
func (b *basicAuthMiddleware) User(ctx *iris2.Context) string {
	return b.config.User(ctx)
}

func (b *basicAuthMiddleware) init() {
	// pass the encoded users from the user's config's Users value
	b.auth = make(encodedUsers, 0, len(b.config.Users))

	for k, v := range b.config.Users {
		fullUser := k + ":" + v
		header := "Basic " + base64.StdEncoding.EncodeToString([]byte(fullUser))
		b.auth = append(b.auth, encodedUser{HeaderValue: header, Username: k, logged: false, expires: DefaultExpireTime})
	}

	// set the auth realm header's value
	b.realmHeaderValue = "Basic realm=" + strconv.Quote(b.config.Realm)

	if b.config.Expires > 0 {
		b.expireEnabled = true
	}
}

func (b *basicAuthMiddleware) findAuth(headerValue string) (auth *encodedUser, found bool) {
	if len(headerValue) == 0 {
		return
	}

	for _, user := range b.auth {
		if user.HeaderValue == headerValue {
			auth = &user
			found = true
			break
		}
	}

	return
}

func (b *basicAuthMiddleware) askForCredentials(ctx *iris2.Context) {
	ctx.SetHeader("WWW-Authenticate", b.realmHeaderValue)
	ctx.SetStatusCode(iris2.StatusUnauthorized)
}

// Serve the actual middleware
func (b *basicAuthMiddleware) Serve(ctx *iris2.Context) {

	if auth, found := b.findAuth(ctx.RequestHeader("Authorization")); !found {
		b.askForCredentials(ctx)
		// don't continue to the next handler
	} else {
		// all ok set the context's value in order to be getable from the next handler
		ctx.Set(b.config.ContextKey, auth.Username)
		if b.expireEnabled {

			if auth.logged == false {
				auth.expires = time.Now().Add(b.config.Expires)
				auth.logged = true
			}

			if time.Now().After(auth.expires) {
				b.askForCredentials(ctx) // ask for authentication again
				return
			}

		}
		ctx.Next() // continue
	}

}
