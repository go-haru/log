package log

import (
	"log"
	"os"

	"github.com/go-haru/field"
)

type Flusher interface {
	Flush() error
}

type Logger interface {
	Debug(v ...any)
	Debugf(format string, v ...any)

	Info(v ...any)
	Infof(format string, v ...any)

	Warn(v ...any)
	Warnf(format string, v ...any)

	Error(v ...any)
	Errorf(format string, v ...any)

	Fatal(v ...any)
	Fatalf(format string, v ...any)

	Panic(v ...any)
	Panicf(format string, v ...any)

	Print(v ...any)
	Printf(format string, v ...any)

	With(v ...field.Field) Logger
	WithName(name string) Logger
	WithLevel(level Level) Logger
	AddDepth(depth int) Logger

	Flusher
	Standard() *log.Logger
}

var logger = NewSimple(os.Stdout, SimpleWithDepth(1))

func Use(l Logger) {
	l = l.AddDepth(1)
	logger = l
}

func Current() Logger {
	var l = logger
	l = l.AddDepth(-1)
	return l
}

func Debug(m ...any) {
	logger.Debug(m...)
}

func Info(m ...any) {
	logger.Info(m...)
}

func Warn(m ...any) {
	logger.Warn(m...)
}

func Error(m ...any) {
	logger.Error(m...)
}

func Panic(m ...any) {
	logger.Panic(m...)
}

func Fatal(m ...any) {
	logger.Fatal(m...)
}

func Debugf(format string, m ...any) {
	logger.Debugf(format, m...)
}

func Infof(format string, m ...any) {
	logger.Infof(format, m...)
}

func Warnf(format string, m ...any) {
	logger.Warnf(format, m...)
}

func Errorf(format string, m ...any) {
	logger.Errorf(format, m...)
}

func Panicf(format string, m ...any) {
	logger.Panicf(format, m...)
}

func Fatalf(format string, m ...any) {
	logger.Fatalf(format, m...)
}

func Printf(format string, m ...any) {
	logger.Infof(format, m...)
}

func With(v ...field.Field) Logger {
	return Current().With(v...)
}

func WithName(name string) Logger {
	return Current().WithName(name)
}
