package logic

import (
	"errors"
	"strconv"
	"time"

	"github.com/antonio-alexander/go-bludgeon/timers/data"
)

const (
	ChangeRateRegistrationLessOrEqualToZero string = "change registration rate less or equal to zero"
	ChangeRateReadLessOrEqualToZero         string = "change read rate less or equal to zero"
	ChangesTimeoutReadLessOrEqualToZero     string = "changes timeout is less or equal to zero"
	ChangesRegistrationIdEmpty              string = "changes registration id empty"
)

const (
	EnvNameChangeRateRegistration string = "BLUDGEON_CHANGE_REGISTRATION_RATE"
	EnvNameChangeRateRead         string = "BLUDGEON_CHANGE_READ_RATE"
	EnvNameChangesTimeout         string = "BLUDGEON_CHANGE_TIMEOUT"
	EnvNameChangesRegistrationId  string = "BLUDGEON_CHANGE_REGISTRATION_ID"
)

const (
	DefaultChangeRateRegistration time.Duration = time.Second
	DefaultChangeRateRead         time.Duration = 10 * time.Second
	DefaultChangesTimeout         time.Duration = 10 * time.Second
)

var (
	DefaultChangesRegistrationId = data.ServiceName
)

var (
	ErrChangeRateRegistrationLessOrEqualToZero = errors.New(ChangeRateRegistrationLessOrEqualToZero)
	ErrChangeRateReadLessOrEqualToZero         = errors.New(ChangeRateReadLessOrEqualToZero)
	ErrChangesTimeoutLessOrEqualToZero         = errors.New(ChangesTimeoutReadLessOrEqualToZero)
	ErrChangesRegistrationIdEmpty              = errors.New(ChangesRegistrationIdEmpty)
)

type Configuration struct {
	ChangeRateRegistration time.Duration `json:"rate_change_registration"`
	ChangeRateRead         time.Duration `json:"rate_change_read"`
	ChangesTimeout         time.Duration `json:"changes_timeout"`
	ChangesRegistrationId  string        `json:"changes_registration_id"`
}

func (c *Configuration) Default() {
	c.ChangeRateRegistration = DefaultChangeRateRegistration
	c.ChangeRateRead = DefaultChangeRateRead
	c.ChangesTimeout = DefaultChangesTimeout
	c.ChangesRegistrationId = DefaultChangesRegistrationId
}

func (c *Configuration) Validate() (err error) {
	if c.ChangeRateRegistration <= 0 {
		return ErrChangeRateRegistrationLessOrEqualToZero
	}
	if c.ChangeRateRead <= 0 {
		return ErrChangeRateReadLessOrEqualToZero
	}
	if c.ChangesTimeout <= 0 {
		return ErrChangesTimeoutLessOrEqualToZero
	}
	if c.ChangesRegistrationId == "" {
		return ErrChangesRegistrationIdEmpty
	}
	return
}

func (c *Configuration) FromEnv(envs map[string]string) {
	if s, ok := envs[EnvNameChangeRateRegistration]; ok && s != "" {
		i, _ := strconv.ParseInt(s, 10, 64)
		c.ChangeRateRegistration = time.Duration(i) * time.Second
	}
	if s, ok := envs[EnvNameChangeRateRead]; ok && s != "" {
		i, _ := strconv.ParseInt(s, 10, 64)
		c.ChangeRateRead = time.Duration(i) * time.Second
	}
	if s, ok := envs[EnvNameChangesTimeout]; ok && s != "" {
		i, _ := strconv.ParseInt(s, 10, 64)
		c.ChangesTimeout = time.Duration(i) * time.Second
	}
	if s, ok := envs[EnvNameChangesRegistrationId]; ok && s != "" {
		c.ChangesRegistrationId = s
	}
}
