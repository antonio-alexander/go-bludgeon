package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Error(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type loggerSimple struct {
	*log.Logger
}

func New(p ...string) interface {
	Logger
} {
	prefix := ""
	if len(p) > 0 {
		prefix = fmt.Sprintf("[%s] ", p[0])
	}
	return &loggerSimple{
		Logger: log.New(os.Stdout, prefix, 0),
	}
}

func (l *loggerSimple) Error(format string, v ...interface{}) {
	l.Printf("Error: "+format, v...)
}

func (l *loggerSimple) Info(format string, v ...interface{}) {
	l.Printf("Info: "+format, v...)
}

func (l *loggerSimple) Debug(format string, v ...interface{}) {
	l.Printf("Debug: "+format, v...)
}
