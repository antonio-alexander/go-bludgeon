package bludgeonconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client/config"
	metajson "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json/config"
	metamysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql/config"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/config"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server/config"

	"github.com/pkg/errors"
)

func Read(configPath, pwd string, envs map[string]string, config interface{}) (err error) {
	var bytes []byte
	var exists bool

	//check if configPath exists, then swtich on the input
	// maintain the pointer for configuration, if config file
	// exists, read from it, otherwise, populate the defaults
	// then attempt to read from the environment and then
	// validate the configuration. In the event the provided
	// configuration is not a supported type, output an error
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		exists = true
	}
	switch c := config.(type) {
	case *Client:
		if !exists {
			DefaultClient(pwd, c)
			if err = fromEnv(pwd, envs, c); err != nil {
				return
			}
		} else {
			if bytes, err = ioutil.ReadFile(configPath); err != nil {
				return
			}
			if err = json.Unmarshal(bytes, &c); err != nil {
				return
			}
		}
		err = Validate(c)
	case *Server:
		if !exists {
			DefaultServer(pwd, c)
			if err = fromEnv(pwd, envs, c); err != nil {
				return
			}
		} else {
			if bytes, err = ioutil.ReadFile(configPath); err != nil {
				return
			}
			if err = json.Unmarshal(bytes, &c); err != nil {
				return
			}
		}
		err = Validate(c)
	default:
		err = errors.Errorf(ErrUnsupportedConfigf, c)
	}

	return
}

func Write(configPath string, config interface{}) (err error) {
	//Switch on the provided configuration and write to disk
	// as needed using the json marshal indent (to make it easier
	// to read)
	switch c := config.(type) {
	case *Client:
		var bytes []byte

		if bytes, err = json.MarshalIndent(&c, "", "    "); err != nil {
			return
		}
		err = ioutil.WriteFile(configPath, bytes, fileMode)
	case *Server:
		var bytes []byte

		if bytes, err = json.MarshalIndent(&c, "", "    "); err != nil {
			return
		}
		err = ioutil.WriteFile(configPath, bytes, fileMode)

	case []byte:
		err = ioutil.WriteFile(configPath, c, fileMode)
	case json.RawMessage:
		err = ioutil.WriteFile(configPath, c, fileMode)
	default:
		err = errors.Errorf(ErrUnsupportedConfigf, c)
	}

	return
}

func fromEnv(pwd string, envs map[string]string, config interface{}) (err error) {
	// var directory string

	//
	switch c := config.(type) {
	case *Client:
		if metaType, ok := envs[EnvNameBludgeonMetaType]; ok {
			c.MetaType = bludgeon.AtoMetaType(metaType)
		}
		if remoteType, ok := envs[EnvNameBludgeonRemoteType]; ok {
			c.RemoteType = bludgeon.AtoRemoteType(remoteType)
		}
		switch c.MetaType {
		case bludgeon.MetaTypeJSON:
			var conf metajson.Configuration

			raw, ok := c.Meta[bludgeon.MetaTypeJSON]
			if !ok {
				err = errors.Errorf(ErrMetaNotFoundf, bludgeon.MetaTypeJSON)

				return
			}
			if err = json.Unmarshal(raw, &conf); err != nil {
				return
			}
			if err = metajson.FromEnv(pwd, envs, &conf); err != nil {
				return
			}
			if c.Meta[bludgeon.MetaTypeJSON], err = json.Marshal(conf); err != nil {
				return
			}
		case bludgeon.MetaTypeMySQL:
			var conf metamysql.Configuration

			raw, ok := c.Meta[bludgeon.MetaTypeMySQL]
			if !ok {
				err = errors.Errorf(ErrMetaNotFoundf, bludgeon.MetaTypeMySQL)

				return
			}
			if err = json.Unmarshal(raw, &conf); err != nil {
				return
			}
			if err = metamysql.FromEnv(pwd, envs, &conf); err != nil {
				return
			}
			if c.Meta[bludgeon.MetaTypeJSON], err = json.Marshal(conf); err != nil {
				return
			}
		}
		switch c.RemoteType {
		case bludgeon.RemoteTypeRest:
			var conf rest.Configuration

			raw, ok := c.Remote[bludgeon.RemoteTypeRest]
			if !ok {
				err = errors.Errorf(ErrRemoteNotFoundf, bludgeon.MetaTypeJSON)

				return
			}
			if err = json.Unmarshal(raw, &conf); err != nil {
				return
			}
			if err = rest.FromEnv(pwd, envs, &conf); err != nil {
				return
			}
			if c.Remote[bludgeon.RemoteTypeRest], err = json.Marshal(conf); err != nil {
				return
			}
		}
		if err = rest.FromEnv(pwd, envs, &c.Rest); err != nil {
			return
		}
		if err = client.FromEnv(pwd, envs, &c.Client); err != nil {
			return
		}
		err = Validate(c)
	case *Server:
		if metaType, ok := envs[EnvNameBludgeonMetaType]; ok {
			c.MetaType = bludgeon.AtoMetaType(metaType)
		}
		switch c.MetaType {
		case bludgeon.MetaTypeJSON:
			var conf metajson.Configuration

			raw, ok := c.Meta[bludgeon.MetaTypeJSON]
			if !ok {
				err = errors.Errorf(ErrMetaNotFoundf, bludgeon.MetaTypeJSON)

				return
			}
			if err = json.Unmarshal(raw, &conf); err != nil {
				return
			}
			if err = metajson.FromEnv(pwd, envs, &conf); err != nil {
				return
			}
			if c.Meta[bludgeon.MetaTypeJSON], err = json.Marshal(conf); err != nil {
				return
			}
		case bludgeon.MetaTypeMySQL:
			var conf metamysql.Configuration

			raw, ok := c.Meta[bludgeon.MetaTypeMySQL]
			if !ok {
				err = errors.Errorf(ErrMetaNotFoundf, bludgeon.MetaTypeMySQL)

				return
			}
			if err = json.Unmarshal(raw, &conf); err != nil {
				return
			}
			if err = metamysql.FromEnv(pwd, envs, &conf); err != nil {
				return
			}
			if c.Meta[bludgeon.MetaTypeJSON], err = json.Marshal(conf); err != nil {
				return
			}
		}
		if err = rest.FromEnv(pwd, envs, &c.Rest); err != nil {
			return
		}
		if err = server.FromEnv(pwd, envs, &c.Server); err != nil {
			return
		}
		err = Validate(c)
	default:
		err = errors.Errorf(ErrUnsupportedConfigf, c)
	}

	return
}

func Validate(config interface{}) (err error) {
	//TODO: this

	return
}
