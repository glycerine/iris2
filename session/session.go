package session

type Session struct {
	storage Storage

	id      string
	data    map[string]interface{}
	flashes map[string][]interface{}
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Save() {
	s.storage.Save(s)
}

func (s *Session) SetFlash(value interface{}, name ...string) {
	nm := "_flash"
	if len(name) > 0 && name[0] != "" {
		nm = name[0]
	}
	s.flashes[nm] = append(s.flashes[nm], value)
}

func (s *Session) Set(key string, value interface{}) {
	s.data[key] = value
}

func (s *Session) GetString(key string) string {
	if s, ok := s.data[key]; ok {
		if str, ok := s.(string); ok {
			return str
		}
	}
	return ""
}

func (s *Session) Delete(key string) {
	if _, ok := s.data[key]; ok {
		delete(s.data, key)
	}
}
