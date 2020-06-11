package bludgeonservercommon

import (
	"fmt"
	"path/filepath"

	json "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
)

func GetConfigFromEnv(pwd string, envs map[string]string) (c Configuration, err error) {
	c = DefaultConfig(pwd)
	//get meta type
	if metaType, ok := envs[EnvNameBludgeonMetaType]; ok {
		c.Meta.Type = metaType
	}
	//get rest address
	if restAddress, ok := envs[EnvNameBludgeonRestAddress]; ok {
		c.Server.Rest.Address = restAddress
	}
	//get rest port
	if restPort, ok := envs[EnvNameBludgeonRestPort]; ok {
		c.Server.Rest.Port = restPort
	}
	//
	for _, metaType := range []string{"json", "mysql"} {
		//get meta configuration
		switch metaType {
		case "json":
			var config json.Configuration

			//get json file
			if jsonFile, ok := envs[EnvNameBludgeonMetaJSONFile]; ok {
				config.File = jsonFile
			}
			c.Meta.Config[metaType] = config
		case "mysql":
			var config mysql.Configuration

			//get database address
			if databaseAddress, ok := envs[EnvNameDatabaseAddress]; ok {
				config.Hostname = databaseAddress
			}
			//get database port
			if databasePort, ok := envs[EnvNameDatabasePort]; ok {
				config.Port = databasePort
			}
			c.Meta.Config[metaType] = config
		default:
			err = fmt.Errorf("Unsupported meta: \"%s\"", metaType)
		}
	}

	return
}

func DefaultConfig(pwd string) (c Configuration) {
	c.Server.Rest.Address = DefaultBludgeonRestAddress
	c.Server.Rest.Port = DefaultBludgeonRestPort
	c.Meta.Type = DefaultBludgeonMeta
	c.Meta.Config = map[string]interface{}{
		"json": json.Configuration{
			File: filepath.Join(pwd, DefaultBludgeonMetaJSONFile),
		},
		"mysql": mysql.Configuration{
			Hostname: DefaultDatabaseAddress,
			Port:     DefaultBludgeonRestPort,
		},
	}

	return
}
