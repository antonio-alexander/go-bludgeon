package client

import "github.com/antonio-alexander/go-bludgeon/timers/meta"

type Client interface {
	meta.TimeSlice
	meta.Timer
}
