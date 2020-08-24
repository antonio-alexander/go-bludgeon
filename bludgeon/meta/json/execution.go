package bludgeonmetajson

import (
	"encoding/json"
	"fmt"

	config "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json/config"
)

func castConfiguration(element interface{}) (c config.Configuration, err error) {
	switch v := element.(type) {
	case json.RawMessage:
		err = json.Unmarshal(v, &c)
	case config.Configuration:
		c = v
	default:
		err = fmt.Errorf("Unsupported type: %t", element)
	}

	return
}
