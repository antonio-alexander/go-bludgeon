package bludgeonmetajson

import (
	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
)

//common constants
const (
	DefaultFile = "data/bludgeon.json"
)

//error constants
const (
	ErrTimerNotFoundf     string = "Timer with id, \"%s\", not found locally"
	ErrTimeSliceNotFoundf string = "TimeSlice with id, \"%s\", not found locally"
	ErrTimerIsArchivedf   string = "Timer with id, \"%s\", is archived"
)

//SerializedData
type SerializedData struct {
	Timers     map[string]bludgeon.Timer     `json:"Timers"`
	TimeSlices map[string]bludgeon.TimeSlice `json:"TimeSlices"`
}

type Configuration struct {
	File string `json:"File"`
}
