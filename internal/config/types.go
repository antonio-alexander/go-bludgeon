package config

import meta "github.com/antonio-alexander/go-bludgeon/internal/meta"

const ErrMetaTypeEmpty string = "meta type empty"

const (
	DefaultConfigFile string    = "bludgeon_service_config.json"
	DefaultConfigPath string    = "config"
	DefaultMetaType   meta.Type = meta.TypeFile
)

const (
	EnvNameMetaType           string = "BLUDGEON_META_TYPE"
	EnvNameServiceRestEnabled string = "BLUDGEON_REST_ENABLED"
	EnvNameServiceGrpcEnabled string = "BLUDGEON_GRPC_ENABLED"
)
