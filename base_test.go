package redis

import (
	"encoding/json"
	"strconv"
	"testing"
)

type People struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func GetPeople() []People {
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

	return peoples
}

func TestSelect(t *testing.T) {
	client, err := NewClient("form")
	if err != nil {
		t.Error(err)
	}
	err = client.Select("menu")
	if err != nil {
		t.Error(err)
	}
}

func TestHash(t *testing.T) {
	name := "test"

	client, err := NewClient("menu", "../../../config/redis.config")
	if err != nil {
		t.Error(err)
	}
	if !client.Exist(name) {

		for index, people := range GetPeople() {
			err = client.HSet(name, strconv.Itoa(index), people)
			if err != nil {
				t.Errorf("HSet Error:	%s", err)
			}
		}

		people, err := client.HGet(name, "0")
		if err != nil {
			t.Errorf("HGet Error:	%s", err)
		} else {
			t.Logf("%s 0:	%s", name, people)
		}

		values, err := client.HValues(name)
		if err != nil {
			t.Errorf("HValues Error:	%s", err)
		} else {
			t.Logf("%s values:	%s", name, values)
		}

		all, err := client.HGetAll(name)
		if err != nil {
			t.Errorf("HGetAll Error:	%s", err)
		} else {
			t.Logf("%s all:	%s", name, all)
		}

		if !client.HExist(name, "1") {
			t.Errorf("HExist Error:	%s 1 does not exist", name)
		}

		keys, err := client.HKeys(name)
		if err != nil {
			t.Errorf("HKeys Error:	%s", err)
		}
		for _, key := range keys {
			err = client.HDel(name, key)
			if err != nil {
				t.Errorf("HDel Error:	%s", err)
			}
		}
	}
}

func TestSet(t *testing.T) {
	name := "test"

	client, err := NewClient("menu", "../../../config/redis.config")
	if err != nil {
		t.Error(err)
	}
	if !client.Exist(name) {

		for _, people := range GetPeople() {
			err = client.SAdd(name, people)
			if err != nil {
				t.Errorf("SAdd Error:	%s", err)
			}
		}

		values, err := client.SValues(name)
		if err != nil {
			t.Errorf("SValues Error:	%s", err)
		} else {
			t.Logf("%s values:	%s", name, values)
		}

		for _, people := range GetPeople() {

			if !client.SExist(name, people) {
				t.Errorf("SExist Error:	%s %v does not exist", name, people)
			}

			err = client.SDel(name, people)
			if err != nil {
				t.Errorf("SDel Error:	%s", err)
			}
		}
	}
}

func TestString(t *testing.T) {
	name := "test"
	people := GetPeople()

	client, err := NewClient("menu", "../../../config/redis.config")
	if err != nil {
		t.Error(err)
	}
	if !client.Exist(name) {

		err = client.Set(name, people[0])
		if err != nil {
			t.Errorf("Set Error:	%s", err)
		}

		value, err := client.Get(name)
		if err != nil {
			t.Errorf("Get Error:	%s", err)
		} else {
			t.Logf("%s value:	%s", name, value)
		}

		if !client.Exist(name) {
			t.Errorf("Exist Error:	%s does not exist", name)
		}

		err = client.Del(name)
		if err != nil {
			t.Errorf("Del Error:	%s", err)
		}

	}
}
