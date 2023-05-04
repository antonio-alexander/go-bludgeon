package rest

import (
	"strconv"
	"time"

	rest "github.com/antonio-alexander/go-bludgeon/internal/rest/client"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// environmental variables
const (
	EnvNameRestAddress  string = "BLUDGEON_TIMERS_REST_ADDRESS"
	EnvNameRestPort     string = "BLUDGEON_TIMERS_REST_PORT"
	EnvNameDisableCache string = "BLUDGEON_TIMERS_DISABLE_CACHE"
)

// defaults
const (
	DefaultPort                   string        = "8012"
	DefaultAddress                string        = "localhost"
	DefaultDisableCache           bool          = false
	DefaultChangesTimeout         time.Duration = 10 * time.Second
	DefaultChangeRateRead         time.Duration = 10 * time.Second
	DefaultChangeRateRegistration time.Duration = 10 * time.Second
)

var DefaultChangesRegistrationId string = uuid.Must(uuid.NewRandom()).String()

type Configuration struct {
	rest.Configuration
	DisableCache           bool          `json:"disable_cache"`
	ChangesRegistrationId  string        `json:"changes_registration_id"`
	ChangesTimeout         time.Duration `json:"changes_timeout"`
	ChangeRateRead         time.Duration `json:"changes_rate_read"`
	ChangeRateRegistration time.Duration `json:"changes_rate_registration"`
}

func (c *Configuration) FromEnvs(envs map[string]string) {
	//Get the address from the environment, then the port
	// then the timeout
	if address, ok := envs[EnvNameRestAddress]; ok {
		c.Address = address
	}
	if port, ok := envs[EnvNameRestPort]; ok {
		c.Port = port
	}
	if timeoutString, ok := envs[rest.EnvNameRestTimeout]; ok {
		if timeoutInt, err := strconv.Atoi(timeoutString); err == nil {
			if timeout := time.Duration(timeoutInt) * time.Second; timeout > 0 {
				c.Timeout = timeout
			}
		}
	}
	if s, ok := envs[EnvNameDisableCache]; ok {
		c.DisableCache, _ = strconv.ParseBool(s)
	}
}

func (r *Configuration) Default() {
	r.Address = DefaultAddress
	r.Port = DefaultPort
	r.Timeout = rest.DefaultTimeout
	r.DisableCache = DefaultDisableCache
	r.ChangesRegistrationId = DefaultChangesRegistrationId
	r.ChangesTimeout = DefaultChangesTimeout
	r.ChangeRateRead = DefaultChangeRateRead
	r.ChangeRateRegistration = DefaultChangeRateRegistration
}

func (c *Configuration) Validate() error {
	//validate that the address isn't empty
	// check if the port is empty, and then ensure
	// that the port is an integer, finally
	// check if the timeout is gt 0
	if c.Address == "" {
		return errors.New(rest.ErrAddressEmpty)
	}
	if c.Port == "" {
		return errors.New(rest.ErrPortEmpty)
	}
	if _, e := strconv.Atoi(c.Port); e != nil {
		return errors.Errorf(rest.ErrPortBadf, c.Port)
	}
	if c.Timeout <= 0 {
		return errors.Errorf(rest.ErrTimeoutBadf, c.Timeout)
	}
	return nil
}
