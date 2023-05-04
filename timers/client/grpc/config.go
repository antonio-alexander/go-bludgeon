package grpc

import (
	"strconv"
	"time"

	grpc "github.com/antonio-alexander/go-bludgeon/internal/grpc/client"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// environmental variables
const (
	EnvNameGrpcAddress  string = "BLUDGEON_TIMERS_GRPC_ADDRESS"
	EnvNameRestPort     string = "BLUDGEON_TIMERS_GRPC_PORT"
	EnvNameDisableCache string = "BLUDGEON_TIMERS_DISABLE_CACHE"
)

// defaults
const (
	DefaultPort                   string        = "8013"
	DefaultAddress                string        = "localhost"
	DefaultDisableCache           bool          = false
	DefaultChangesTimeout         time.Duration = 10 * time.Second
	DefaultChangeRateRead         time.Duration = 10 * time.Second
	DefaultChangeRateRegistration time.Duration = 10 * time.Second
)

var DefaultChangesRegistrationId string = uuid.Must(uuid.NewRandom()).String()

type Configuration struct {
	grpc.Configuration
	DisableCache           bool          `json:"disable_cache"`
	ChangesRegistrationId  string        `json:"changes_registration_id"`
	ChangesTimeout         time.Duration `json:"changes_timeout"`
	ChangeRateRead         time.Duration `json:"changes_rate_read"`
	ChangeRateRegistration time.Duration `json:"changes_rate_registration"`
}

func (r *Configuration) Default() {
	r.Address = DefaultAddress
	r.Port = DefaultPort
	r.DisableCache = DefaultDisableCache
	r.ChangesRegistrationId = DefaultChangesRegistrationId
	r.ChangesTimeout = DefaultChangesTimeout
	r.ChangeRateRead = DefaultChangeRateRead
	r.ChangeRateRegistration = DefaultChangeRateRegistration
}

func (c *Configuration) FromEnvs(envs map[string]string) {
	//Get the address from the environment, then the port
	// then the timeout
	if address, ok := envs[EnvNameGrpcAddress]; ok {
		c.Address = address
	}
	if port, ok := envs[EnvNameRestPort]; ok {
		c.Port = port
	}
}

func (c *Configuration) Validate() error {
	//validate that the address isn't empty
	// check if the port is empty, and then ensure
	// that the port is an integer, finally
	// check if the timeout is gt 0
	if c.Address == "" {
		return errors.New(grpc.ErrAddressEmpty)
	}
	if c.Port == "" {
		return errors.New(grpc.ErrPortEmpty)
	}
	if _, e := strconv.Atoi(c.Port); e != nil {
		return errors.Errorf(grpc.ErrPortBadf, c.Port)
	}
	return nil
}
