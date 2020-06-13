package bludgeonconfigserver

import (
	"encoding/json"
	"time"

	mjson "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql"
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
	// c.Rest, _ = json.Marshal(rest.Configuration{
	// 	Address: "127.0.0.1",
	// 	Port:    "8080",
	// 	Timeout: 10 * time.Second,
	// })

	return
}

// func DefaultConfig(pwd string) (c Configuration) {
// 	c.Server.Rest.Address = DefaultBludgeonRestAddress
// 	c.Server.Rest.Port = DefaultBludgeonRestPort
// 	c.Meta.Type = DefaultBludgeonMeta
// 	c.Meta.Config = map[string]interface{}{
// 		"json": json.Configuration{
// 			File: filepath.Join(pwd, DefaultBludgeonMetaJSONFile),
// 		},
// 		"mysql": mysql.Configuration{
// 			Hostname: DefaultDatabaseAddress,
// 			Port:     DefaultBludgeonRestPort,
// 		},
// 	}

// 	return
// }
