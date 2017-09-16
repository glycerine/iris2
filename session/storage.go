package session

type Storage interface {
	Init(*Handler)
	Load(string) *Session
	Exist(string) bool
	New(string) *Session
	Save(*Session)
}
