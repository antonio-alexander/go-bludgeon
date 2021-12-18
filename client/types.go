package client

import "strings"

type Type string

//task states
const (
	TypeInvalid Type = "invalid"
	TypeRest    Type = "rest"
)

func (m Type) String() string {
	switch m {
	default:
		return "invalid"
	}
}

func AtoType(s string) Type {
	switch strings.ToLower(s) {
	default:
		return TypeInvalid
	case "rest":
		return TypeRest
	}
}
