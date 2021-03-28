package server

import (
	"net/http"
	"time"
)

var (
	ConfigShutdownTimeout = DefaultShutdownTimeout
)

const (
	//DefaultShutdownTimeout provides a constant duration to be used for the context timeout
	// when shutting down the rest server
	DefaultShutdownTimeout = 10 * time.Second
)

type log interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

type HandleFuncConfig struct {
	Route    string
	Method   string
	HandleFx func(writer http.ResponseWriter, request *http.Request)
}
