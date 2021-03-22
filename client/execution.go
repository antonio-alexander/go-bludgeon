package client

import (
	"os"
	"os/user"
	"path/filepath"
)

const (
	DefaultFolder            = ".bludgeon"
	DefaultConfigurationFile = "config/bludgeon_config.json"
	DefaultCacheFile         = "data/bludgeon_cache.json"
)

func Files(pwd string, config interface{}) (configFile, cacheFile string, err error) {
	switch config.(type) {
	case *Configuration, Configuration:
		//get the user from os
		if u, _ := user.Current(); u != nil {
			//get the home dir
			if homedir := u.HomeDir; homedir != "" {
				bludgeonDirectory := filepath.Join(homedir, DefaultFolder)
				//TODO: attempt to create default folder
				if _, err = os.Stat(bludgeonDirectory); os.IsNotExist(err) {
					//create folder
					if err = os.MkdirAll(bludgeonDirectory, 0700); err != nil {
						return
					}
				}
				//generate all the paths
				configFile = filepath.Join(bludgeonDirectory, DefaultConfigurationFile)
				cacheFile = filepath.Join(bludgeonDirectory, DefaultCacheFile)
			}
		}
	}
	//check if the config file is empty
	if configFile == "" {
		configFile = filepath.Join(pwd, DefaultConfigurationFile)
	}
	//check if the cache file is empty
	if cacheFile == "" {
		cacheFile = filepath.Join(pwd, DefaultCacheFile)
	}

	return
}
