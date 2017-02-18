package file

import (
	"gopkg.in/vmihailenco/msgpack.v2"
	"io/ioutil"
	"os"
)

// Database structure for the file-storage
type Database struct {
	path string
}

// New returns a new session storage instance
func New(p string) *Database {
	return &Database{path: p}
}

// Load loads the values to the underline
func (d *Database) Load(sid string) map[string]interface{} {
	values := make(map[string]interface{})

	val, err := ioutil.ReadFile(d.path + "/" + sid)
	if err == nil {
		err = msgpack.Unmarshal(val, &values)
		if err != nil {
			println("Filestorage deserialize error: " + err.Error())
		}
	}

	return values

}

// serialize the values to be stored as strings inside the session storage
func serialize(values map[string]interface{}) []byte {
	val, err := msgpack.Marshal(values)
	if err != nil {
		println("Filestorage serialize error: " + err.Error())
	}

	return val
}

// Update updates the session storage
func (d *Database) Update(sid string, newValues map[string]interface{}) {
	if len(newValues) == 0 {
		go os.Remove(d.path + "/" + sid)
	} else {
		ioutil.WriteFile(d.path+"/"+sid, serialize(newValues), 0600)
	}

}
