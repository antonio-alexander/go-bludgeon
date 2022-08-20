package internal

import (
	"os"
)

// Main is used to init an instance of the bludgeon changes service
func Main(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) error {
	config := getConfig(pwd, envs)
	logger, parameters := parameterize(config)
	if err := configure(pwd, envs, append(parameters, logger)...); err != nil {
		return err
	}
	logger.Info("version \"%s\" (%s@%s)", Version, GitBranch, GitCommit)
	if err := initialize(parameters...); err != nil {
		return err
	}
	defer func() { shutdown(parameters...) }()
	//wait for signal, then shutdown
	<-chSignalInt
	return nil
}
