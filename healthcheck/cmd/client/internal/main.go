package internal

import (
	"context"
	"time"
)

// Main is used to start an instance of the bludgeon healthcheck service
func Main(pwd string, args []string, envs map[string]string) error {
	config := getConfig(pwd, args, envs)
	logger, client := parameterize(config)
	if err := configure(pwd, envs, []interface{}{logger, client}...); err != nil {
		return err
	}
	logger.Info("version \"%s\" (%s@%s)", Version, GitBranch, GitCommit)
	if err := initialize([]interface{}{logger, client}...); err != nil {
		return err
	}
	defer shutdown([]interface{}{logger, client}...)
	healthcheck, err := client.HealthCheck(context.Background())
	if err != nil {
		return err
	}
	logger.Trace("healthcheck: %v", time.Unix(0, healthcheck.Time))
	return nil
}
