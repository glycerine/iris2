package sessions

import (
	"sync"
	"time"

	"github.com/go-iris2/iris2"
)

type (
	// provider contains the sessions and external databases (load and update).
	// It's the session memory manager
	provider struct {
		// we don't use RWMutex because all actions have read and write at the same action function.
		// (or write to a *session's value which is race if we don't lock)
		// narrow locks are fasters but are useless here.
		mu       sync.Mutex
		sessions map[string]*session
		database Database
	}
)

// newProvider returns a new sessions provider
func newProvider(db Database) *provider {
	return &provider{
		sessions: make(map[string]*session, 0),
		database: db,
	}
}

// newSession returns a new session from sessionid
func (p *provider) newSession(sid string, expires time.Duration) *session {
	sess := &session{
		sid:      sid,
		provider: p,
		values:   p.loadSessionValues(sid),
		flashes:  make(map[string]*flashMessage),
	}

	if expires > 0 { // if not unlimited life duration and no -1 (cookie remove action is based on browser's session)
		sess.timeout = time.AfterFunc(expires, func() {
			p.Destroy(sid)
		})
	}

	return sess
}

func (p *provider) loadSessionValues(sid string) map[string]interface{} {
	if p.database != nil {
		dbValues, err := p.database.Load(sid)
		if dbValues != nil && err == nil {
			return dbValues
		}
	}

	return make(map[string]interface{})
}

// Init creates the session  and returns it
func (p *provider) Init(sid string, expires time.Duration) iris2.Session {
	newSession := p.newSession(sid, expires)
	p.mu.Lock()
	p.sessions[sid] = newSession
	p.mu.Unlock()
	return newSession
}

// Read returns the store which sid parameter belongs
func (p *provider) Read(sid string, expires time.Duration) iris2.Session {
	p.mu.Lock()
	sess, found := p.sessions[sid]
	p.mu.Unlock()
	if found {
		if p.sessions[sid].timeout != nil {
			p.sessions[sid].timeout.Reset(expires)
		}
		sess.runFlashGC()
		return sess
	}

	// When it is not in p.sessions it must be in de database
	// p.Init loads it from there
	return p.Init(sid, expires)
}

func (p *provider) Exist(sid string) bool {
	p.mu.Lock()
	_, found := p.sessions[sid]
	p.mu.Unlock()
	if !found && p.database != nil {
		_, err := p.database.Load(sid)
		if err == nil {
			found = true
		}
	}
	return found
}

// Destroy destroys the session, removes all sessions and flash values,
// the session itself and updates the registered session databases,
// this called from sessionManager which removes the client's cookie also.
func (p *provider) Destroy(sid string) {
	p.mu.Lock()
	if sess, found := p.sessions[sid]; found {
		sess.values = nil
		sess.flashes = nil
		delete(p.sessions, sid)
		p.updateDb(sid, nil)
	}
	p.mu.Unlock()
}

// DestroyAll removes all sessions
// from the server-side memory (and database if registered).
// Client's session cookie will still exist but it will be reseted on the next request.
func (p *provider) DestroyAll() {
	p.mu.Lock()
	for _, sess := range p.sessions {
		delete(p.sessions, sess.ID())
		p.updateDb(sess.ID(), nil)
	}
	p.mu.Unlock()
}

func (p *provider) updateDb(sid string, values map[string]interface{}) {
	if p.database != nil {
		p.database.Update(sid, values)
	}
}
