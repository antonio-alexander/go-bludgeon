package bludgeonmetajsonconfig

//environmental variables
const (
	EnvNameBludgeonMetaJSONFile string = "BLUDGEON_META_JSON_FILE"
)

//defaults
const (
	DefaultBludgeonMetaJSONFile string = "data/bludgeon.json"
)

type Configuration struct {
	File string
}
