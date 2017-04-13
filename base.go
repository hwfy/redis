//Copyright (c) 2017, hwfy

//Permission to use, copy, modify, and/or distribute this software for any
//purpose with or without fee is hereby granted, provided that the above
//copyright notice and this permission notice appear in all copies.

//THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
//WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
//MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
//ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
//WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
//ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
//OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package redis

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type (
	// redis connection pool
	client struct {
		pool        *redis.Pool
		DataSources map[string]dataSource `json:"dataSources"`
	}
	// redis dataSource
	dataSource struct {
		Num   int    `json:"num"`
		Addr  string `json:"addr"`
		Port  string `json:"port"`
		Pass  string `json:"pass"`
		Pool  int    `json:"pool"`
		Cache string `json:"cache"`
	}
)

func (c *client) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	rc := c.pool.Get()
	defer rc.Close()

	return rc.Do(commandName, args...)
}

// marshal serialized into json
func marshal(value interface{}) interface{} {
	actual := reflect.ValueOf(value)
	switch actual.Kind() {
	case reflect.Map, reflect.Struct, reflect.Ptr, reflect.Array, reflect.Slice:
		if actual.Type().String() != "[]uint8" {
			value, _ = json.Marshal(value)
		}
	}
	return value
}

// unmarshal if the value of the map is json is deserialized
func unmarshal(value interface{}) ([]byte, error) {
	switch actual := value.(type) {
	case map[string]string:
		result := make(map[string]interface{}, len(actual))
		for k, v := range actual {
			var object interface{}

			if err := json.Unmarshal([]byte(v), &object); err != nil {
				object = v
			}
			result[k] = object
		}
		return json.Marshal(result)

	case [][]byte:
		result := make([]interface{}, len(actual))
		for i, v := range actual {
			var object interface{}

			if err := json.Unmarshal(v, &object); err != nil {
				object = string(v)
			}
			result[i] = object
		}
		return json.Marshal(result)
	}
	return json.Marshal(value)
}

// SetExpire set the expiration time of key, time is seconds
// if the key does not exist to return an error
func (c *client) SetExpire(key string, time int) error {
	exsit, err := redis.Bool(c.do("EXPIRE", key, time))
	if err != nil {
		return err
	}
	if !exsit {
		return errors.New(key + " does not exist in the database")
	}
	return nil
}

// Select switch the database, name is the dataSource in the configuration file
func (c *client) Select(name string) error {
	if c == nil {
		return errors.New("initialize redis failed")
	}
	db, has := c.DataSources[name]
	if !has {
		return errors.New("the configuration file is missing node " + name)
	}
	_, err := c.do("SELECT", db.Num)

	return err
}

// Close close the redis connection pool
func (c *client) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}

// NewClient init redis, name is the dataSource in the configuration file
// can specify the configuration path, if addr is empty, it is set to 127.0.0.1
// port is empty, it is set to 6379; cache is empty, it is set to 60s.
func NewClient(name string, path ...string) (*client, error) {
	cfg := "../config/redis.config"

	if path != nil && path[0] != "" {
		cfg = path[0]
	}
	bytes, err := ioutil.ReadFile(cfg)
	if err != nil {
		return nil, errors.New("read configuration file failed, " + err.Error())
	}
	redis := new(client)
	if err = json.Unmarshal(bytes, redis); err != nil {
		return nil, err
	}
	db, has := redis.DataSources[name]
	if !has {
		return nil, errors.New("the configuration file is missing node " + name)
	}
	switch "" {
	case db.Addr:
		db.Addr = "127.0.0.1"
	case db.Port:
		db.Port = "6379"
	case db.Cache:
		db.Cache = "60s"
	}
	if !strings.HasSuffix(db.Cache, "s") {
		db.Cache += "s"
	}
	addr := db.Addr + ":" + db.Port

	timeOut, _ := time.ParseDuration(db.Cache)

	redis.pool = newPool(addr, db.Pass, db.Num, db.Pool, timeOut)

	return redis, nil
}

func newPool(server, password string, dbNum, poolSize int, cacheTime time.Duration) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     poolSize,
		IdleTimeout: cacheTime,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			_, selecterr := c.Do("SELECT", dbNum)
			if selecterr != nil {
				c.Close()
				return nil, selecterr
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
