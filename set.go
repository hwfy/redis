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
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// SExist determine whether the element is in the collection
func (c *client) SExist(key string, value interface{}) bool {
	exist, err := redis.Bool(c.do("SISMEMBER", key, marshal(value)))
	if err != nil {
		return false
	}
	return exist
}

// SAdd add an element to the collection
// if the boolean type is stored 0 or 1
func (c *client) SAdd(key string, value interface{}) error {
	if value == nil {
		return fmt.Errorf("invalid value nil")
	}
	_, err := c.do("SADD", key, marshal(value))

	return err
}

// SValues gets all the members in the collection
// if the key does not exist will return an error
func (c *client) SValues(key string) ([]byte, error) {
	values, err := redis.ByteSlices(c.do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("%s does not exist in the set database", key)
	}
	return unmarshal(values)
}

// SDel removes the specified element in the collection
func (c *client) SDel(key string, value interface{}) error {
	_, err := c.do("SREM", key, marshal(value))

	return err
}
