package config

import (
	"encoding/json"
	"errors"
)

func Get(item interface{}, configKey string, config any) error {
	if configKey == "" {
		return nil
	}
	bytes, err := json.Marshal(item)
	if err != nil {
		return nil
	}
	items := make(map[string]interface{})
	err = json.Unmarshal(bytes, &items)
	if err != nil {
		return nil
	}
	for key, item := range items {
		if key != configKey {
			continue
		}
		bytes, err := json.Marshal(item)
		if err != nil {
			return err
		}
		return json.Unmarshal(bytes, config)
	}
	return errors.New(ErrConfigurationNotFound)
}
