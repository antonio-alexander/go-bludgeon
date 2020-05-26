package bludgeonclient

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	json "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	mysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/server/api"
)

func initClientFolders(pwd string) (configFile, cacheFile, jsonFile string, err error) {
	//get the user from os
	if u, _ := user.Current(); u != nil {
		//get the home dir
		if homedir := u.HomeDir; homedir != "" {
			bludgeonDirectory := filepath.Join(homedir, bludgeon.DefaultFolder)
			//TODO: attempt to create default folder
			if _, err = os.Stat(bludgeonDirectory); os.IsNotExist(err) {
				//create folder
				if err = os.MkdirAll(bludgeonDirectory, 0700); err != nil {
					return
				}
			}
			//generate all the paths
			configFile = filepath.Join(bludgeonDirectory, bludgeon.DefaultConfigurationFile)
			cacheFile = filepath.Join(bludgeonDirectory, bludgeon.DefaultCacheFile)
			jsonFile = filepath.Join(bludgeonDirectory, json.DefaultFile)
		}
	}
	//check if the config file is empty
	if configFile == "" {
		configFile = filepath.Join(pwd, bludgeon.DefaultConfigurationFile)
	}
	//check if the cache file is empty
	if cacheFile == "" {
		cacheFile = filepath.Join(pwd, bludgeon.DefaultCacheFile)
	}
	//check if the cache file is empty
	if jsonFile == "" {
		jsonFile = filepath.Join(pwd, json.DefaultFile)
	}

	return
}

func initMeta(metaType string, config interface{}) (meta interface {
	bludgeon.MetaOwner
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
}, err error) {
	// filepath.Join(pwd, "bludgeon.json")
	switch s := strings.ToLower(metaType); s {
	case "json":
		//create metajson
		m := json.NewMetaJSON()
		//initialize metajson
		if err = m.Initialize(config); err != nil {
			return
		}
		meta = m
	case "mysql":
		m := mysql.NewMetaMySQL()
		//connect
		if err = m.Initialize(config); err != nil {
			return
		}
		meta = m
	default:
		err = fmt.Errorf("meta type unsupported: %s", s)
	}

	return
}

func initRemote(remoteType string, config interface{}) (remote interface {
	bludgeon.RemoteOwner
	bludgeon.RemoteTimer
	bludgeon.RemoteTimeSlice
}, err error) {
	//switch on rest type
	switch remoteType {
	case "rest":
		//create rest remote
		r := rest.NewRemote()
		if err = r.Initialize(config); err != nil {
			return
		}
		remote = r
	default:
		err = fmt.Errorf("remote type unsupported: %s", remoteType)
	}

	return
}
