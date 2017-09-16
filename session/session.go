package session

type Session struct {
	storage Storage

	id   string
	data map[string]interface{}
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Save() {
	s.storage.Save(s)
}
