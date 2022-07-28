package logger

import (
	"fmt"
	"log"
	"os"
)

type logger struct {
	*log.Logger
	config *Configuration
}

func New(parameters ...interface{}) interface {
	Logger
} {
	l := &logger{}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *Configuration:
			l.config = p
		}
	}
	if l.config == nil {
		l.config = &Configuration{}
		l.config.Default()
	}
	l.Logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", l.config.Prefix), 0)
	l.Info("logger configured with level \"%s\"", l.config.Level)
	return l
}

func (l *logger) Error(format string, v ...interface{}) {
	if l.config.Level >= Error {
		l.Printf("Error: "+format, v...)
	}
}

func (l *logger) Info(format string, v ...interface{}) {
	if l.config.Level >= Info {
		l.Printf("Info: "+format, v...)
	}
}

func (l *logger) Debug(format string, v ...interface{}) {
	if l.config.Level >= Debug {
		l.Printf("Debug: "+format, v...)
	}
}

func (l *logger) Trace(format string, v ...interface{}) {
	if l.config.Level >= Trace {
		l.Printf("Trace: "+format, v...)
	}
}
