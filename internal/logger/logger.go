package logger

import (
	"log"
	"os"

	"github.com/antonio-alexander/go-bludgeon/common"
)

type logger struct {
	*log.Logger
	prefix string
}

func New(prefix string) common.Logger {
	return &logger{
		Logger: log.New(os.Stdout, "", 0),
		prefix: prefix,
	}
}

//Error
func (l *logger) Error(err error, v ...interface{}) {
	if l.prefix != "" {
		l.Printf("[%s] Error: %s", l.prefix, err)
	} else {
		l.Printf("Error: %s", l.prefix, err)
	}
}

//Info
func (l *logger) Info(format string, v ...interface{}) {
	if l.prefix != "" {
		l.Printf("["+l.prefix+"] Info: "+format, v...)
	} else {
		l.Printf("Info: "+format, v...)
	}
}

//Debug
func (l *logger) Debug(format string, v ...interface{}) {
	if l.prefix != "" {
		l.Printf("["+l.prefix+"] Debug: "+format, v...)
	} else {
		l.Printf("Debug: "+format, v...)
	}
}
