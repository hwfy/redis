# redis
A redis client that can access json data, currently supports hash tables, string, set

# Installation
> go get github.com/hwfy/redis

# Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hwfy/redis"
)

type People struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	name := "test"
	data := `
	[
		{
			"name": "bill",
			"age": 64
		},
		{
			"name": "hwfy",
			"age": 26
		}
	]
	`
	var peoples []People
	json.Unmarshal([]byte(data), &peoples)

	//menu is the datasource in the configuration file
	//you can also specify a path: 
	//redis.NewClient("menu","../../../config/redis.config")

	client, err := redis.NewClient("menu")
	if err != nil {
		panic(err)
	}
	for index, people := range peoples {
		err = client.HSet(name, strconv.Itoa(index), people)
		if err != nil {
			panic(err)
		}
	}
	keys, err := client.HKeys(name)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s fields: %s\n", name, keys)

	people, err := client.HGet(name, "0")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s 0: %s\n", name, people)

	values, err := client.HValues(name)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s values: %s\n", name, values)

	all, err := client.HGetAll(name)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s all: %s\n", name, all)

	// OutPut:
	// test fields: [0 1]
	// test 0: {"name":"bill","age":64}
	// test values: [{"age":64,"name":"bill"},{"age":26,"name":"hwfy"}]
	// test all: {"0":{"age":64,"name":"bill"},"1":{"age":26,"name":"hwfy"}}
}
```
