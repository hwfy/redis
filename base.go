// Copyright 2017 hwfy
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

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

//marshal serialized into json
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

	err = json.Unmarshal(bytes, redis)
	if err != nil {
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
