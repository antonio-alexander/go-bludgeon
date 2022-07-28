package internal

import (
	"os"
)

//Main is used to start an instance of the bludgeon timers service
func Main(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) error {
	config := getConfig(pwd, envs)
	logger := getLogger(config)
	logger.Info("version \"%s\" (%s@%s)", Version, GitBranch, GitCommit)
	meta, err := startMeta(config, logger)
	if err != nil {
		return err
	}
	defer meta.Shutdown()
	logic, err := startLogic(config, logger, meta)
	if err != nil {
		return err
	}
	stopServices, err := startServices(config, logic, logger)
	if err != nil {
		return err
	}
	defer stopServices()
	<-chSignalInt
	return nil
}
