package client

import (
	"encoding/json"
	"strings"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/client"
	cli "github.com/antonio-alexander/go-bludgeon/client/cli"
	rest "github.com/antonio-alexander/go-bludgeon/client/rest"
	data "github.com/antonio-alexander/go-bludgeon/data"
	logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	logger_simple "github.com/antonio-alexander/go-bludgeon/internal/logger/simple"
	logic "github.com/antonio-alexander/go-bludgeon/internal/logic"

	"github.com/pkg/errors"
)

func execute(client logic.Logic, options cli.Options, logger logger.Logger) (err error) {
	switch options.ObjectType {
	default:
		return errors.Errorf("unsupported object: %s", options.ObjectType)
	case data.ObjectTypeTimer:
		var timer data.Timer

		switch strings.ToLower(options.Command) {
		default:
			return errors.Errorf("unsupported command: %s", options.Command)
		case "create":
			if timer, err = client.TimerCreate(); err != nil {
				return
			}
		case "read":
			if timer, err = client.TimerRead(options.Timer.UUID); err != nil {
				return
			}
		case "update":
			if timer, err = client.TimerUpdate(options.Timer); err != nil {
				return
			}
		case "delete":
			return client.TimerDelete(options.Timer.UUID)
		case "start":
			if timer, err = client.TimerStart(options.Timer.UUID, time.Unix(0, options.Timer.Start)); err != nil {
				return
			}
		case "pause":
			//REVIEW: i don't think this works
			if timer, err = client.TimerPause(options.Timer.UUID, time.Unix(0, options.Timer.Finish)); err != nil {
				return
			}
		case "submit":
			if timer, err = client.TimerSubmit(options.Timer.UUID, time.Unix(0, options.Timer.Finish)); err != nil {
				return
			}
		}
		bytes, _ := json.MarshalIndent(&timer, "", " ")
		logger.Info("timer %s", string(bytes))
	}

	return
}

func getClient(config *Configuration) (logic.Logic, error) {
	switch config.Client.Type {
	default:
		return nil, errors.Errorf("Unsupported client: %s", config.Client.Type)
	case client.TypeRest:
		return rest.New(*config.Client.Rest), nil
	}
}

func Main(pwd string, args []string, envs map[string]string) error {
	logger := logger_simple.New("bludgeon-client")
	logger.Info("version %s (%s@%s)", Version, GitBranch, GitCommit)
	options, err := cli.Parse(pwd, args, envs)
	if err != nil {
		return err
	}
	configPath, cacheFile, err := Files(pwd)
	if err != nil {
		return err
	}
	cache := &cache{}
	if err := cache.Read(cacheFile); err != nil {
		return err
	}
	if options.Timer.UUID == "" {
		options.Timer.UUID = cache.TimerID
	}
	config := NewConfiguration()
	if err = config.Read(configPath, pwd, envs); err != nil {
		return err
	}
	client, err := getClient(config)
	if err != nil {
		return err
	}
	return execute(client, options, logger)
}
