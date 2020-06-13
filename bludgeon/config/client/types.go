package bludgeonconfigclient

import (
	"encoding/json"

	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client/common"
)

type Configuration struct {
	Meta struct {
		Type   string                     `json:"Meta"`
		Config map[string]json.RawMessage `json:"Config"`
	}
	Remote struct {
		Type   string                     `json:"Type"`
		Config map[string]json.RawMessage `json:"Config"`
	} `json:"Remote"`
	Client client.Configuration `json:"Client"`
	Rest   struct {
		Address string `json:"Address"`
		Port    string `json:"Port"`
	} `json:"Rest"`
}
