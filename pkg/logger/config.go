package logger

const (
	EnvNameLevel  string = "BLUDGEON_LOG_LEVEL"
	EnvNamePrefix string = "BLUDGEON_LOG_PREFIX"
)

const (
	DefaultLevel Level = Info
)

type Configuration struct {
	Level  Level  `json:"level"`
	Prefix string `json:"prefix"`
}

func (c *Configuration) Default() {
	c.Level = DefaultLevel
}

func (c *Configuration) FromEnv(envs map[string]string) {
	if level, ok := envs[EnvNameLevel]; ok {
		c.Level = AtoLogLevel(level)
	}
	if prefix, ok := envs[EnvNamePrefix]; ok {
		c.Prefix = prefix
	}
}
