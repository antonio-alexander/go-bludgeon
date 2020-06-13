package bludgeonconfigclient

import (
	"encoding/json"
	"time"

	mjson "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/server/api"
)

func configDefault(jsonFile string) (c Configuration) {
	//meta
	c.Meta.Type = "json"
	c.Meta.Config = make(map[string]json.RawMessage)
	c.Meta.Config["json"], _ = json.Marshal(mjson.Configuration{
		File: jsonFile,
	})
	c.Meta.Config["mysql"], _ = json.Marshal(mysql.Configuration{
		Hostname:        "127.0.0.1",
		Port:            "3306",
		Username:        "bludgeon",
		Password:        "bludgeon",
		Database:        "bludgeon",
		ParseTime:       false,
		UseTransactions: true,
		Timeout:         10 * time.Second,
		// DataSource:      "",
		// FilePath:        "",
	})
	c.Remote.Type = "rest"
	c.Remote.Config = make(map[string]json.RawMessage)
	c.Remote.Config["rest"], _ = json.Marshal(rest.Configuration{
		Address: "127.0.0.1",
		Port:    "8080",
		Timeout: 10 * time.Second,
	})

	return
}
