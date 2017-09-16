package session

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type RamStorage struct {
	c *cache.Cache
	h *Handler
}

func (r *RamStorage) Init(h *Handler) {
	r.h = h
	r.c = cache.New(time.Duration(h.expireSeconds)*time.Second, 10*time.Minute)
}

func (r *RamStorage) Load(sessID string) *Session {
	data, ok := r.c.Get(sessID)
	if ok {
		return &Session{
			storage: r,
			id:      sessID,
			data:    data.(map[string]interface{}),
		}
	}
	return nil
}

func (r *RamStorage) Exist(sessID string) bool {
	_, ok := r.c.Get(sessID)
	return ok
}

func (r *RamStorage) New(sessID string) *Session {
	s := &Session{
		data:    make(map[string]interface{}),
		id:      sessID,
		storage: r,
	}
	r.c.Add(sessID, s.data, cache.DefaultExpiration)
	return s
}

func (r *RamStorage) Save(s *Session) {
	r.c.Set(s.id, s.data, cache.DefaultExpiration)
}
