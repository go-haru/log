package log

import "strings"

type Level uint

const (
	DebugLevel Level = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarningLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "info"
	}
}

func ParseLevel(str string) (Level, bool) {
	switch strings.ToLower(str) {
	case "debug":
		return DebugLevel, true
	case "info":
		return InfoLevel, true
	case "warning":
		return WarningLevel, true
	case "error":
		return ErrorLevel, true
	case "fatal":
		return FatalLevel, true
	default:
		return InfoLevel, false
	}
}
