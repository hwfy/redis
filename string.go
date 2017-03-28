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
	"errors"

	"github.com/garyburd/redigo/redis"
)

// Exist determine whether the key exists
func (c *client) Exist(key string) bool {
	exist, err := redis.Bool(c.do("EXISTS", key))
	if err != nil {
		return false
	}
	return exist
}

// Set set the string value
// if the boolean type is stored 0 or 1
func (c *client) Set(key string, value interface{}) error {
	if value == nil {
		return errors.New("invalid value nil")
	}
	_, err := c.do("SET", key, marshal(value))

	return err
}

// Get get the string value
func (c *client) Get(key string) ([]byte, error) {
	return redis.Bytes(c.do("GET", key))
}

// Del delete the specified key
// if the key does not exist will return an error
func (c *client) Del(key string) error {
	_, err := c.do("DEL", key)

	return err
}

// Keys get all the keys in the database
func (c *client) Keys() ([]string, error) {
	return redis.Strings(c.do("KEYS", "*"))
}
