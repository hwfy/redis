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
	//If the value of the slice is json is deserialized
	var object interface{}
	result := make([]interface{}, len(values))

	for i, v := range values {
		err = json.Unmarshal(v, &object)
		if err != nil {
			object = string(v)
		}
		result[i] = object
	}
	return json.Marshal(result)
}

// SDel removes the specified element in the collection
func (c *client) SDel(key string, value interface{}) error {
	_, err := c.do("SREM", key, marshal(value))

	return err
}
