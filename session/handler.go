package session

import (
	"time"

	"github.com/pborman/uuid"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	cookieName    string
	expireSeconds uint
	storage       Storage
}

func New() *Handler {
	h := &Handler{
		cookieName:    "SessionID",
		expireSeconds: 86400,
		storage:       &RamStorage{},
	}
	h.storage.Init(h)
	return h
}

func (h *Handler) Config() {
}

func (h *Handler) Start(ctx *fasthttp.RequestCtx) *Session {
	c := ctx.Request.Header.Cookie(h.cookieName)
	ip := ctx.RemoteIP().String()
	if h.storage.Exist(ip+"_"+string(c)) && len(c) > 0 {
		h.updateSessionCookie(ctx, ip+"_"+string(c))
		return h.storage.Load(ip + "_" + string(c))
	}
	for {
		sess := uuid.NewRandom().String()
		if !h.storage.Exist(ip + "_" + sess) {
			h.updateSessionCookie(ctx, ip+"_"+string(c))
			return h.storage.New(ip + "_" + sess)
		}
	}
}

func (h *Handler) updateSessionCookie(ctx *fasthttp.RequestCtx, sessionName string) {
	c := &fasthttp.Cookie{}
	c.SetKey(h.cookieName)
	c.SetValue(sessionName)
	c.SetPath("/")
	c.SetHTTPOnly(true)
	c.SetExpire(time.Now().Add(time.Second * time.Duration(h.expireSeconds)))
	ctx.Response.Header.SetCookie(c)
}
