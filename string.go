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
