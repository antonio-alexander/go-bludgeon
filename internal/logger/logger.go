package logger

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/config"
)

type logger struct {
	*log.Logger
	config *Configuration
}

func New() interface {
	Logger
	Printer
	internal.Configurer
} {
	config := new(Configuration)
	config.Default()
	return &logger{config: new(Configuration)}
}

func (l *logger) Configure(items ...interface{}) error {
	var c *Configuration

	for _, item := range items {
		switch v := item.(type) {
		case config.Envs:
			c = new(Configuration)
			c.FromEnv(v)
		case *Configuration:
			c = v
		}
	}
	if c == nil {
		return errors.New(config.ErrConfigurationNotFound)
	}
	l.config = c
	l.Logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", l.config.Prefix), log.Ltime|log.Ldate|log.Lmsgprefix)
	l.Info("logger configured with level \"%s\"", l.config.Level)
	return nil
}

func (l *logger) Error(format string, v ...interface{}) {
	if l.config.Level >= Error {
		l.Printf("[error] "+format, v...)
	}
}

func (l *logger) Info(format string, v ...interface{}) {
	if l.config.Level >= Info {
		l.Printf("[info] "+format, v...)
	}
}

func (l *logger) Debug(format string, v ...interface{}) {
	if l.config.Level >= Debug {
		l.Printf("[debug] "+format, v...)
	}
}

func (l *logger) Trace(format string, v ...interface{}) {
	if l.config.Level >= Trace {
		l.Printf("[trace] "+format, v...)
	}
}

type nullLogger struct{}

func NewNullLogger(parameters ...interface{}) interface {
	Logger
	Printer
} {
	return &nullLogger{}
}

func (n *nullLogger) Error(format string, v ...interface{})  {}
func (n *nullLogger) Info(format string, v ...interface{})   {}
func (n *nullLogger) Debug(format string, v ...interface{})  {}
func (n *nullLogger) Trace(format string, v ...interface{})  {}
func (n *nullLogger) Print(v ...interface{})                 {}
func (n *nullLogger) Printf(format string, v ...interface{}) {}
func (n *nullLogger) Println(v ...interface{})               {}
