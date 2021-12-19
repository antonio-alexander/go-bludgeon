package logger

type Logger interface {
	Error(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}
