package config

const ErrMetaTypeEmpty string = "meta type empty"

const (
	DefaultConfigFile string = "bludgeon.json"
	DefaultConfigPath string = "config"
)

type (
	Envs map[string]string
	Args []string
)
