package service

import "strings"

type Type string

//type constants
const (
	TypeInvalid Type = "invalid"
	TypeREST    Type = "rest"
)

func (m Type) String() string {
	switch m {
	default:
		return "invalid"
	case TypeREST:
		return "rest"
	}
}

func AtoType(s string) Type {
	switch strings.ToLower(s) {
	default:
		return TypeInvalid
	case "rest":
		return TypeREST
	}
}
