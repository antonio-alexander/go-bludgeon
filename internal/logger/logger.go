package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/antonio-alexander/go-bludgeon/data"
)

type logger struct {
	*log.Logger
}

func New(p ...string) interface {
	data.Logger
} {
	prefix := ""
	if len(p) > 0 {
		prefix = fmt.Sprintf("[%s] ", p[0])
	}
	return &logger{
		Logger: log.New(os.Stdout, prefix, 0),
	}
}

func (l *logger) Error(format string, v ...interface{}) {
	l.Printf("Error: "+format, v...)
}

func (l *logger) Info(format string, v ...interface{}) {
	l.Printf("Info: "+format, v...)
}

func (l *logger) Debug(format string, v ...interface{}) {
	l.Printf("Debug: "+format, v...)
}
