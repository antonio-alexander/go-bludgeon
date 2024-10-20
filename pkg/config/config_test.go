package config_test

import (
	"testing"

	"github.com/antonio-alexander/go-bludgeon/pkg/config"

	"github.com/stretchr/testify/assert"
)

type configNested struct {
	One int `json:"one"`
	Two int `json:"Two"`
}

type masterConfig struct {
	configNested `json:"configNested"`
	Three        int `json:"three"`
}

func TestConfig(t *testing.T) {
	var c masterConfig
	var cc configNested

	err := config.Get(c, "configNested", &cc)
	assert.Nil(t, err)
	assert.IsType(t, c, &configNested{})
}
