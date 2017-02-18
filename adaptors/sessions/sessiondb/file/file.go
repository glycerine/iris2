package file

import (
	"fmt"
	"github.com/go-iris2/iris2/adaptors/sessions"
	"gopkg.in/vmihailenco/msgpack.v2"
	"io/ioutil"
	"os"
)

// fileStorage structure for the file-storage
type fileStorage struct {
	path string
}

// New returns a new session storage instance
func New(p string) sessions.Database {
	return &fileStorage{path: p}
}

// Load loads the values to the underline
func (d *fileStorage) Load(sid string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	val, err := ioutil.ReadFile(d.path + "/" + sid)
	if err != nil {
		return nil, fmt.Errorf("could not read session: %v", err)
	}

	err = msgpack.Unmarshal(val, &values)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize session: %v", err)
	}

	return values, nil

}

// Update updates the session storage
func (d *fileStorage) Update(sid string, newValues map[string]interface{}) {
	if newValues == nil || len(newValues) == 0 {
		go os.Remove(d.path + "/" + sid)
	} else {
		val, err := msgpack.Marshal(newValues)
		if err == nil {
			ioutil.WriteFile(d.path+"/"+sid, val, 0600)
		}
	}
}
