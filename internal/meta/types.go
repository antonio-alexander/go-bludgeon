package meta

import (
	"strings"

	"github.com/pkg/errors"
)

//Type is a special string that can be used to identify the type
// of meta
type Type string

const (
	TypeInvalid Type = "invalid"
	TypeMemory  Type = "memory"
	TypeFile    Type = "file"
	TypeMySQL   Type = "mysql"
)

//String can be used to convert a meta type into a string
// type
func (m Type) String() string {
	switch m {
	case TypeMemory:
		return "memory"
	case TypeFile:
		return "file"
	case TypeMySQL:
		return "mysql"
	default:
		return "invalid"
	}
}

//AtoType can be used to attempt to convert a string to a
// meta type
func AtoType(s string) Type {
	switch strings.ToLower(s) {
	case "memory":
		return TypeMemory
	case "file":
		return TypeFile
	case "mysql":
		return TypeMySQL
	default:
		return TypeInvalid
	}
}

func ErrUnsupportedMeta(metaType Type) error {
	return errors.Errorf("unsupported meta: %s", metaType)
}

//Owner contains methods that shuold only be used by the constructor of the pointer
// and allows use of functions that can affect the underlying pointer
type Owner interface {
	//Shutdown can be used to "stop" the meta and any asynchronous
	// processes
	Shutdown()
}
