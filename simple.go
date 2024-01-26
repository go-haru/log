package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-haru/field"
)

type SimpleLoggerOption func(*simpleLogger)

func SimpleWithName(name string) SimpleLoggerOption {
	return func(l *simpleLogger) {
		l.name = name
		l.logger.SetFlags(l.logger.Flags() | log.Lmsgprefix)
		l.logger.SetPrefix("[" + name + "]")
	}
}

func SimpleWithData(fields ...field.Field) SimpleLoggerOption {
	return func(l *simpleLogger) {
		l.dataArr = append(l.dataArr, fields...)
		l.buildDataStr()
	}
}

func SimpleWithLevel(level Level) SimpleLoggerOption {
	return func(l *simpleLogger) { l.level = level }
}

func SimpleWithDepth(depth int) SimpleLoggerOption {
	return func(l *simpleLogger) { l.depth = depth }
}

func SimpleAddDepth(depth int) SimpleLoggerOption {
	return func(l *simpleLogger) { l.depth += depth }
}

func NewSimple(output io.Writer, opts ...SimpleLoggerOption) Logger {
	var sysLogger = log.New(
		output, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile,
	)
	_logger := &simpleLogger{logger: sysLogger}
	for _, opt := range opts {
		if opt != nil {
			opt(_logger)
		}
	}
	return _logger
}

type simpleLogger struct {
	name    string
	dataArr field.Fields
	dataStr string
	depth   int
	level   Level
	logger  *log.Logger
}

func (s *simpleLogger) Standard() *log.Logger { return s.logger }

func (s *simpleLogger) Debug(v ...any) {
	s.errHandler(s.logger.Output(s.depth+2, "[D] "+fmt.Sprint(v...)+s.dataStr))
}

func (s *simpleLogger) Debugf(format string, v ...any) {
	s.errHandler(s.logger.Output(s.depth+2, "[D] "+fmt.Sprintf(format, v...)+s.dataStr))
}

func (s *simpleLogger) Info(v ...any) {
	s.errHandler(s.logger.Output(s.depth+2, "[I] "+fmt.Sprint(v...)+s.dataStr))
}

func (s *simpleLogger) Infof(format string, v ...any) {
	s.errHandler(s.logger.Output(s.depth+2, "[I] "+fmt.Sprintf(format, v...)+s.dataStr))
}

func (s *simpleLogger) Warn(v ...any) {
	s.errHandler(s.logger.Output(s.depth+2, "[W] "+fmt.Sprint(v...)+s.dataStr))
}

func (s *simpleLogger) Warnf(format string, v ...any) {
	s.errHandler(s.logger.Output(s.depth+2, "[W] "+fmt.Sprintf(format, v...)+s.dataStr))
}

func (s *simpleLogger) Error(v ...any) {
	s.errHandler(s.logger.Output(s.depth+2, "[E] "+fmt.Sprint(v...)+s.dataStr))
}

func (s *simpleLogger) Errorf(format string, v ...any) {
	s.errHandler(s.logger.Output(s.depth+2, "[E] "+fmt.Sprintf(format, v...)+s.dataStr))
}

func (s *simpleLogger) Fatal(v ...any) {
	content := fmt.Sprint(v...)
	s.errHandler(s.logger.Output(s.depth+2, "[F] "+content+s.dataStr))
	_ = s.Flush()
	os.Exit(-1)
}

func (s *simpleLogger) Fatalf(format string, v ...any) {
	content := fmt.Sprintf(format, v...)
	s.errHandler(s.logger.Output(s.depth+2, "[F] "+content+s.dataStr))
	_ = s.Flush()
	os.Exit(-1)
}

func (s *simpleLogger) Panic(v ...any) {
	content := fmt.Sprint(v...)
	s.errHandler(s.logger.Output(s.depth+2, "[P] "+content+s.dataStr))
	_ = s.Flush()
	panic(content)
}

func (s *simpleLogger) Panicf(format string, v ...any) {
	content := fmt.Sprintf(format, v...)
	s.errHandler(s.logger.Output(s.depth+2, "[P] "+content+s.dataStr))
	_ = s.Flush()
	panic(content)
}

func (s *simpleLogger) Print(v ...any) {
	s.errHandler(s.logger.Output(s.depth+2, s.printPrefix()+fmt.Sprint(v...)+s.dataStr))
}

func (s *simpleLogger) Printf(format string, v ...any) {
	s.errHandler(s.logger.Output(s.depth+2, s.printPrefix()+fmt.Sprintf(format, v...)+s.dataStr))
}

func (s *simpleLogger) printPrefix() string {
	switch s.level {
	case DebugLevel:
		return "[D] "
	case InfoLevel:
		return "[I] "
	case WarningLevel:
		return "[W] "
	case ErrorLevel:
		return "[E] "
	case FatalLevel:
		return "[F] "
	default:
		return "[I] "
	}
}

func (s *simpleLogger) clone() *simpleLogger {
	return &simpleLogger{
		name:    s.name,
		dataArr: s.dataArr,
		dataStr: s.dataStr,
		depth:   s.depth,
		level:   s.level,
		logger:  s.logger,
	}
}

func (s *simpleLogger) With(v ...field.Field) Logger {
	var l = s.clone()
	SimpleWithData(v...)(l)
	return l
}

func (s *simpleLogger) WithName(name string) Logger {
	var l = s.clone()
	SimpleWithName(name)(l)
	return l
}

func (s *simpleLogger) WithLevel(level Level) Logger {
	var l = s.clone()
	SimpleWithLevel(level)(l)
	return l
}

func (s *simpleLogger) AddDepth(depth int) Logger {
	var l = s.clone()
	SimpleAddDepth(depth)(l)
	return l
}

func (s *simpleLogger) buildDataStr() {
	if len(s.dataArr) == 0 {
		s.dataStr = ""
		return
	}
	var builder bytes.Buffer
	builder.Write([]byte{0x20, '#', 0x20})
	s.errHandler(s.dataArr.EncodeJSON(&builder))
	s.dataStr = builder.String()
}

func (s *simpleLogger) Flush() error {
	if _flusher, ok := s.logger.Writer().(Flusher); ok {
		return _flusher.Flush()
	}
	return nil
}

func (s *simpleLogger) errHandler(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "log output fail: %v\n", err)
	}
}
