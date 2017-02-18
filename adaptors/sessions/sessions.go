// Package sessions for Iris2, original by kataras, updated by rikvdh
// Based on kataras/go-sessions.
//
package sessions

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-iris2/iris2"
	"math/rand"
)

type (
	// Sessions is the start point of this package
	// contains all the registered sessions and manages them
	Sessions interface {
		// Adapt is used to adapt this sessions manager as an iris2.SessionsPolicy
		// to an Iris station.
		// It's being used by the framework, developers should not actually call this function.
		Adapt(*iris2.Policies)

		// Start starts the session for the particular net/http request
		Start(http.ResponseWriter, *http.Request) iris2.Session

		// Destroy kills the net/http session and remove the associated cookie
		Destroy(http.ResponseWriter, *http.Request)

		// DestroyByID removes the session entry
		// from the server-side memory (and database if registered).
		// Client's session cookie will still exist but it will be reseted on the next request.
		//
		// It's safe to use it even if you are not sure if a session with that id exists.
		DestroyByID(string)
		// DestroyAll removes all sessions
		// from the server-side memory (and database if registered).
		// Client's session cookie will still exist but it will be reseted on the next request.
		DestroyAll()
	}

	// sessions contains the cookie's name, the provider and a duration for GC and cookie life expire
	sessions struct {
		config   Config
		provider *provider
	}
)

// New returns a new fast, feature-rich sessions manager
// it can be adapted to an Iris station
func New(cfg Config) Sessions {
	return &sessions{
		config:   cfg.Validate(),
		provider: newProvider(cfg.SessionStorage),
	}
}

func (s *sessions) Adapt(frame *iris2.Policies) {
	// for newcomers this maybe looks strange:
	// Each policy is an adaptor too, so they all can contain an Adapt.
	// If they contains an Adapt func then the policy is an adaptor too and this Adapt func is called
	// by Iris on .Adapt(...)
	policy := iris2.SessionsPolicy{
		Start:   s.Start,
		Destroy: s.Destroy,
	}

	policy.Adapt(frame)
}

// Start starts the session for the particular net/http request
func (s *sessions) Start(res http.ResponseWriter, req *http.Request) iris2.Session {
	var sess iris2.Session

	clientIP := req.RemoteAddr
	idx := strings.IndexByte(clientIP, ',')
	if idx > 0 {
		clientIP = clientIP[0:idx]
	}

	cookieName := GetCookie(s.config.Cookie, req)
	sessionID := clientIP + "_" + cookieName
	if cookieName != "" && s.provider.Exist(sessionID) {
		sess = s.provider.Read(sessionID, s.config.Expires)
	} else {
		for {
			cookieName = sessionIDGenerator(s.config.CookieLength)
			sessionID = clientIP + "_" + cookieName
			if !s.provider.Exist(sessionID) {
				break
			}
		}
		sess = s.provider.Init(sessionID, s.config.Expires)
	}
	// We always use AddCookie
	SetCookie(s.buildCookie(cookieName, req.URL.Host), res)

	return sess
}

// Destroy kills the net/http session and remove the associated cookie
func (s *sessions) Destroy(res http.ResponseWriter, req *http.Request) {
	cookieName := GetCookie(s.config.Cookie, req)
	if cookieName == "" { // nothing to destroy
		return
	}
	RemoveCookie(s.config.Cookie, res, req)
	clientIP := req.RemoteAddr
	idx := strings.IndexByte(clientIP, ',')
	if idx > 0 {
		clientIP = clientIP[0:idx]
	}
	s.provider.Destroy(clientIP + "_" + cookieName)
}

// DestroyByID removes the session entry
// from the server-side memory (and database if registered).
// Client's session cookie will still exist but it will be reseted on the next request.
//
// It's safe to use it even if you are not sure if a session with that id exists.
// Works for both net/http
func (s *sessions) DestroyByID(sid string) {
	s.provider.Destroy(sid)
}

// DestroyAll removes all sessions
// from the server-side memory (and database if registered).
// Client's session cookie will still exist but it will be reseted on the next request.
// Works for both net/http
func (s *sessions) DestroyAll() {
	s.provider.DestroyAll()
}

func (s *sessions) buildCookie(sid, host string) *http.Cookie {
	cookie := http.Cookie{
		Name:     s.config.Cookie,
		Value:    sid,
		Path:     "/",
		HttpOnly: s.config.HTTPOnly,
	}

	if !s.config.DisableSubdomainPersistence {
		requestDomain := host
		if portIdx := strings.IndexByte(requestDomain, ':'); portIdx > 0 {
			requestDomain = requestDomain[0:portIdx]
		}
		if IsValidCookieDomain(requestDomain) {
			// RFC2109, we allow level 1 subdomains, but no further
			// if we have localhost.com , we want the localhost.cos.
			// so if we have something like: mysubdomain.localhost.com we want the localhost here
			// if we have mysubsubdomain.mysubdomain.localhost.com we want the .mysubdomain.localhost.com here
			// slow things here, especially the 'replace' but this is a good and understable( I hope) way to get the be able to set cookies from subdomains & domain with 1-level limit
			if dotIdx := strings.LastIndexByte(requestDomain, '.'); dotIdx > 0 {
				// is mysubdomain.localhost.com || mysubsubdomain.mysubdomain.localhost.com
				s := requestDomain[0:dotIdx] // set mysubdomain.localhost || mysubsubdomain.mysubdomain.localhost
				if secondDotIdx := strings.LastIndexByte(s, '.'); secondDotIdx > 0 {
					//is mysubdomain.localhost ||  mysubsubdomain.mysubdomain.localhost
					s = s[secondDotIdx+1:] // set to localhost || mysubdomain.localhost
				}
				// replace the s with the requestDomain before the domain's siffux
				subdomainSuff := strings.LastIndexByte(requestDomain, '.')
				if subdomainSuff > len(s) { // if it is actual exists as subdomain suffix
					requestDomain = strings.Replace(requestDomain, requestDomain[0:subdomainSuff], s, 1) // set to localhost.com || mysubdomain.localhost.com
				}
			}
			// finally set the .localhost.com (for(1-level) || .mysubdomain.localhost.com (for 2-level subdomain allow)
			cookie.Domain = "." + requestDomain // . to allow persistence
		}

	}

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	if s.config.Expires >= 0 {
		if s.config.Expires == 0 { // unlimited life
			cookie.Expires = CookieExpireUnlimited
		} else { // > 0
			cookie.Expires = time.Now().Add(s.config.Expires)
		}
		cookie.MaxAge = int(cookie.Expires.Sub(time.Now()).Seconds())
	}
	return &cookie
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// SessionIDGenerator generates a random string of size n
func sessionIDGenerator(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
