package bludgeonmetajsonconfig

//error constants
const (
	ErrFileEmpty string = "File is empty"
)

//environmental variables
const (
	EnvNameFile string = "BLUDGEON_META_JSON_FILE"
)

//defaults
const (
	DefaultFile string = "data/bludgeon.json"
)

type Configuration struct {
	File string
}
