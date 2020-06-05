package bludgeonclient

import (
	"encoding/json"
	"io/ioutil"
)

const (
	CacheFile string = "bludgeon_cache_client.json"
)

type Cache struct {
	TimerID string `json:"TimerID"`
}

func CacheRead(file string) (cache Cache, err error) {
	var bytes []byte

	if bytes, err = ioutil.ReadFile(file); err != nil {
		return
	}
	err = json.Unmarshal(bytes, &cache)

	return
}

func CacheWrite(file string, cache Cache) (err error) {
	var bytes []byte

	if bytes, err = json.Marshal(&cache); err != nil {
		return
	}
	err = ioutil.WriteFile(file, bytes, 0644)

	return
}
