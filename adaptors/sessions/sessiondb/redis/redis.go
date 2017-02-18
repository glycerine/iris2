package redis

import (
	"fmt"
	"github.com/go-iris2/iris2/adaptors/sessions"
	"github.com/go-iris2/iris2/adaptors/sessions/sessiondb/redis/service"
	"gopkg.in/vmihailenco/msgpack.v2"
)

// redisStorage the redis redisStorage for q sessions
type redisStorage struct {
	redis *service.Service
}

// New returns a new redis redisStorage
func New(cfg ...service.Config) sessions.Database {
	return &redisStorage{redis: service.New(cfg...)}
}

// Config returns the configuration for the redis server bridge, you can change them
func (d *redisStorage) Config() *service.Config {
	return d.redis.Config
}

// Load loads the values to the underline
func (d *redisStorage) Load(sid string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	if !d.redis.Connected {
		d.redis.Connect()
		_, err := d.redis.PingPong()
		if err != nil {
			return nil, fmt.Errorf("no connection to redis: %v", err)
		}
	}

	val, err := d.redis.GetBytes(sid)
	if err != nil {
		return nil, fmt.Errorf("fetching session from redis failed: %v", err)
	}

	if err := msgpack.Unmarshal(val, &values); err != nil {
		return nil, fmt.Errorf("error decoding session: %v", err)
	}

	return values, nil
}

// Update updates the real redis store
func (d *redisStorage) Update(sid string, newValues map[string]interface{}) {
	if len(newValues) == 0 {
		d.redis.Delete(sid)
	} else {
		data, err := msgpack.Marshal(newValues)
		if err == nil {
			d.redis.Set(sid, data) //set/update all the values
		}
	}
}
