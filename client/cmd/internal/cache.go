package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type cache struct {
	TimerID string `json:"TimerID"`
}

func (c *cache) Read(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		c.Default()
		if err = c.Write(file); err != nil {
			return err
		}
		return nil
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	cache := &cache{}
	if err = json.Unmarshal(bytes, &cache); err != nil {
		return err
	}
	c.TimerID = cache.TimerID
	return nil
}

func (c *cache) Write(file string) error {
	bytes, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, bytes, 0644)
}

func (c *cache) Default() {
	c.TimerID = ""
}
