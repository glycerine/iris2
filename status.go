package iris2

import (
	"net/http"
	"sync"
)

// ErrorHandlers contains all custom http errors.
// A custom http error handler is just a handler with its status code.
type ErrorHandlers struct {
	// Handlers the map which actually contains the errors.
	// Use the declared functions to get, set or fire an error.
	handlers map[int]Handler
	mu       sync.RWMutex
}

// Register registers a handler to a http status
func (e *ErrorHandlers) Register(statusCode int, handler Handler) {
	e.mu.Lock()
	if e.handlers == nil {
		e.handlers = make(map[int]Handler)
	}
	func(statusCode int, handler Handler) {
		e.handlers[statusCode] = HandlerFunc(func(ctx *Context) {
			if w, ok := ctx.IsRecording(); ok {
				w.Reset()
			}
			ctx.SetStatusCode(statusCode)
			handler.Serve(ctx)
		})
	}(statusCode, handler)
	e.mu.Unlock()
}

// Get returns the handler which is responsible for
// this 'statusCode' http error.
func (e *ErrorHandlers) Get(statusCode int) Handler {
	e.mu.RLock()
	h := e.handlers[statusCode]
	e.mu.RUnlock()
	if h == nil {
		return nil
	}
	return h
}

// GetOrRegister trys to return the handler which is responsible
// for the 'statusCode', if it was nil then it creates
// a new one, registers that to the error list and returns that.
func (e *ErrorHandlers) GetOrRegister(statusCode int) Handler {
	h := e.Get(statusCode)
	if h != nil {
		return h
	}
	// create a new one
	h = HandlerFunc(func(ctx *Context) {
		if w, ok := ctx.IsRecording(); ok {
			w.Reset()
		}
		ctx.SetStatusCode(statusCode)
		if _, err := ctx.WriteString(http.StatusText(statusCode)); err != nil {
			ctx.Log("error from a pre-defined error handler while trying to send an http error: %s",
				err.Error())
		}
	})
	e.mu.Lock()
	e.handlers[statusCode] = h
	e.mu.Unlock()
	return h
}

// Fire fires an error based on the `statusCode`
func (e *ErrorHandlers) Fire(statusCode int, ctx *Context) {
	h := e.GetOrRegister(statusCode)
	h.Serve(ctx)
}
