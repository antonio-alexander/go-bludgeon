package config

const (
	ErrMetaTypeEmpty         string = "meta type empty"
	ErrConfigurationNotFound string = "no configuration found"
)

const (
	DefaultConfigFile string = "bludgeon.json"
	DefaultConfigPath string = "config"
)

type (
	Envs map[string]string
	Args []string
)
