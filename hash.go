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

var hashNil = " does not exist in the hash table database"

// HGetAll gets all the fields and values of the hash table
// if the key does not exist will return an error
func (c *client) HGetAll(key string) ([]byte, error) {
	maps, err := redis.StringMap(c.do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	if len(maps) == 0 {
		return nil, errors.New(key + hashNil)
	}
	return unmarshal(maps)
}

// HValues gets the value of all fields in the hash table
// if the key does not exist will return an error
func (c *client) HValues(key string) ([]byte, error) {
	values, err := redis.ByteSlices(c.do("HVALS", key))
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, errors.New(key + hashNil)
	}
	return unmarshal(values)
}

// HKeys gets all the fields in the hash table
// if the key does not exist will return an error
func (c *client) HKeys(key string) ([]string, error) {
	keys, err := redis.Strings(c.do("HKEYS", key))
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, errors.New(key + hashNil)
	}
	return keys, nil
}

// HGet gets the value of the given field in the hash table
func (c *client) HGet(key, field string) ([]byte, error) {
	return redis.Bytes(c.do("HGET", key, field))
}

// HDel delete the specified field in the hash table
func (c *client) HDel(key, field string) error {
	_, err := c.do("HDEL", key, field)

	return err
}

// HEXISTS determine if the field in the hash table exist
func (c *client) HExist(key string, field interface{}) bool {
	exist, err := redis.Bool(c.do("HEXISTS", key, field))
	if err != nil {
		return false
	}
	return exist
}

// HPut sets the value of the field in the hash table
// if the boolean type is stored 0 or 1
func (c *client) HSet(key, field string, value interface{}) error {
	if value == nil {
		return errors.New("invalid value nil")
	}
	_, err := c.do("HSET", key, field, marshal(value))

	return err
}

// HLen get the total number of fields in the hash table key +1
// if the key not exist to return 0, the error returns -1
func (c *client) HLen(key string) int {
	num, err := redis.Int(c.do("HLEN", key))
	if err != nil {
		return -1
	}
	for {
		if c.HExist(key, num) {
			num++
		} else {
			break
		}
	}
	return num
}
