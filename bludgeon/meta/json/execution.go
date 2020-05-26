package bludgeonmetajson

import (
	"encoding/json"
	"fmt"
)

func castConfiguration(element interface{}) (c Configuration, err error) {

	switch v := element.(type) {
	case json.RawMessage:
		err = json.Unmarshal(v, &c)
	case Configuration:
		c = v
	default:
		err = fmt.Errorf("Unsupported type: %t", element)
	}

	return
}
