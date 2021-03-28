package logger

import (
	"log"
	"os"

	"github.com/antonio-alexander/go-bludgeon/common"
)

type logger struct {
	*log.Logger
}

func New() common.Logger {
	return &logger{
		Logger: log.New(os.Stdout, "", 0),
	}
}

//Error
func (l *logger) Error(err error, v ...interface{}) {
	l.Printf("Error: %s", err)
}

//Info
func (l *logger) Info(format string, v ...interface{}) {
	l.Printf("Info: "+format, v...)
}

//Debug
func (l *logger) Debug(format string, v ...interface{}) {
	l.Printf("Debug: "+format, v...)
}
