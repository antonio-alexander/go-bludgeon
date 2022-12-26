package logger

import "strings"

type Logger interface {
	Error(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Trace(format string, v ...interface{})
}
type Printer interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type Level int

const (
	Error Level = 1
	Info  Level = 2
	Debug Level = 3
	Trace Level = 4
)

func (l Level) String() string {
	switch l {
	default:
		return ""
	case Error:
		return "error"
	case Info:
		return "info"
	case Debug:
		return "debug"
	case Trace:
		return "trace"
	}
}

func AtoLogLevel(a string) Level {
	switch strings.ToLower(a) {
	default:
		return Error
	case "info":
		return Info
	case "debug":
		return Debug
	case "trace":
		return Trace
	}
}
