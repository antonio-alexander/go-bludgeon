package restserver

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	ErrAddressEmpty        string = "address is empty"
	ErrPortEmpty           string = "port is empty"
	ErrPortBadf            string = "port is a non-integer: %s"
	ErrShutdownTimeoutBadf string = "shutdown timeout is lte to 0: %v"
)

const (
	EnvNameAddress          string = "BLUDGEON_REST_ADDRESS"
	EnvNamePort             string = "BLUDGEON_REST_PORT"
	EnvNameShutdownTimeout  string = "BLUDGEON_REST_SHUTDOWN_TIMEOUT"
	EnvNameAllowCredentials string = "BLUDGEON_ALLOW_CREDENTIALS"
	EnvNameAllowedOrigins   string = "BLUDGEON_ALLOWED_ORIGINS"
	EnvNameAllowedMethods   string = "BLUDGEON_ALLOWED_METHODS"
	EnvNameCorsDebug        string = "BLUDGEON_CORS_DEBUG"
	EnvNameCorsDisabled     string = "BLUDGEON_DISABLE_CORS"
)

const (
	DefaultAddress          string        = "127.0.0.1"
	DefaultPort             string        = "8080"
	DefaultShutdownTimeout  time.Duration = 10 * time.Second
	DefaultAllowCredentials bool          = true
	DefaultCorsDebug        bool          = false
)

type Configuration struct {
	Address          string        `json:"address"`
	Port             string        `json:"port"`
	Timeout          time.Duration `json:"timeout"`
	ShutdownTimeout  time.Duration `json:"shutdown_timeout"`
	AllowedOrigins   []string      `json:"allowed_origins"`
	AllowedMethods   []string      `json:"allowed_methods"`
	AllowCredentials bool          `json:"allow_credentials"`
	CorsDebug        bool          `json:"cors_debug"`
	CorsDisabled     bool          `json:"cors_disabled"`
}

func (c *Configuration) FromEnv(envs map[string]string) {
	//Get the address from the environment, then the port
	// then the timeout
	if address, ok := envs[EnvNameAddress]; ok {
		c.Address = address
	}
	if port, ok := envs[EnvNamePort]; ok {
		c.Port = port
	}
	if shutdownTimeoutString, ok := envs[EnvNameShutdownTimeout]; ok {
		if shutdownTimeoutInt, err := strconv.Atoi(shutdownTimeoutString); err == nil {
			if timeout := time.Duration(shutdownTimeoutInt) * time.Second; timeout > 0 {
				c.ShutdownTimeout = timeout
			}
		}
	}
	if allowCredentialsString, ok := envs[EnvNameAllowCredentials]; ok {
		if allowCredentials, err := strconv.ParseBool(allowCredentialsString); err == nil {
			c.AllowCredentials = allowCredentials
		}
	}
	if allowedOrigins, ok := envs[EnvNameAllowedOrigins]; ok {
		c.AllowedOrigins = strings.Split(allowedOrigins, ",")
	}
	if allowedMethods, ok := envs[EnvNameAllowedMethods]; ok {
		c.AllowedMethods = strings.Split(allowedMethods, ",")
	}
	if corsDebugString, ok := envs[EnvNameCorsDebug]; ok {
		if corsDebug, err := strconv.ParseBool(corsDebugString); err == nil {
			c.CorsDebug = corsDebug
		}
	}
	if corsDisabledString, ok := envs[EnvNameCorsDisabled]; ok {
		if corsDisabled, err := strconv.ParseBool(corsDisabledString); err == nil {
			c.CorsDisabled = corsDisabled
		}
	}
}

func (c *Configuration) Validate() (err error) {
	//validate that the address isn't empty
	// check if the port is empty, and then ensure
	// that the port is an integer, finally
	// check if the timeout is lte 0
	if c.Port == "" {
		return errors.New(ErrPortEmpty)
	}
	if _, e := strconv.Atoi(c.Port); e != nil {
		return errors.Errorf(ErrPortBadf, c.Port)
	}
	if c.ShutdownTimeout <= 0 {
		return errors.Errorf(ErrShutdownTimeoutBadf, c.ShutdownTimeout)
	}
	return
}
