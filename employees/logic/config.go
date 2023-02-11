package logic

import (
	"errors"
	"strconv"
	"time"
)

const (
	FrequencyChangeRegistrationLessOrEqualToZero string = "change registration frequency less or equal to zero"
	ChangesTimeoutReadLessOrEqualToZero          string = "changes timeout is less or equal to zero"
)

const (
	EnvNameFrequencyChangeRegistration string = "BLUDGEON_CHANGE_FREQUENCY_REGISTRATION"
	EnvNameChangesTimeout              string = "BLUDGEON_CHANGE_TIMEOUT"
)

const (
	DefaultFrequencyChangeRegistration time.Duration = time.Second
	DefaultChangesTimeout              time.Duration = 10 * time.Second
)

var (
	ErrFrequencyChangeRegistrationLessOrEqualToZero = errors.New(FrequencyChangeRegistrationLessOrEqualToZero)
	ErrChangesTimeoutLessOrEqualToZero              = errors.New(ChangesTimeoutReadLessOrEqualToZero)
)

type Configuration struct {
	FrequencyChangeRegistration time.Duration `json:"frequency_change_registration"`
	ChangesTimeout              time.Duration `json:"changes_timeout"`
}

func (c *Configuration) Default() {
	c.FrequencyChangeRegistration = DefaultFrequencyChangeRegistration
	c.ChangesTimeout = DefaultChangesTimeout
}

func (c *Configuration) Validate() (err error) {
	if c.FrequencyChangeRegistration <= 0 {
		return ErrFrequencyChangeRegistrationLessOrEqualToZero
	}
	if c.ChangesTimeout <= 0 {
		return ErrChangesTimeoutLessOrEqualToZero
	}
	return
}

func (c *Configuration) FromEnv(envs map[string]string) {
	if s, ok := envs[EnvNameFrequencyChangeRegistration]; ok && s != "" {
		i, _ := strconv.ParseInt(s, 10, 64)
		c.FrequencyChangeRegistration = time.Duration(i) * time.Second
	}
	if s, ok := envs[EnvNameChangesTimeout]; ok && s != "" {
		i, _ := strconv.ParseInt(s, 10, 64)
		c.ChangesTimeout = time.Duration(i) * time.Second
	}
}
