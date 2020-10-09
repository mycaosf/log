package log

import (
	"errors"
)

//param:
// - LogTypeNull, LogTypeStdout: ignore
// - LogTypeSyslog: tag
// - LogTypeFile: fileName
func New(logType, minLevel int, param string) (l Log, err error) {
	switch logType {
	case LogTypeNull:
		l = NewNull()
	case LogTypeStdout:
		l = NewStdout(minLevel)
	case LogTypeStdoutColor:
		l = NewStdoutColor(minLevel)
	case LogTypeFile:
		l, err = NewFile(param, minLevel)
	case LogTypeSyslog:
		l, err = NewFile(param, minLevel)
	default:
		err = ErrInvalidLogType
	}

	return
}

func (p *logBase) Fatalf(format string, args ...interface{}) {
}

func (p *logBase) Errorf(format string, args ...interface{}) {
}

func (p *logBase) Infof(format string, args ...interface{}) {
}

func (p *logBase) Debugf(format string, args ...interface{}) {
}

func (p *logBase) SetMinLevel(level int) {
	p.level = level
}

func (p *logBase) MinLevel() int {
	return p.level
}

func (p *logBase) Write(b []byte) (n int, err error) {
	n = len(b)

	return
}

func (p *logBase) Close() error {
	return nil
}

func NewNull() Log {
	return &logBase{level: LogLevelFatal}
}

type logBase struct {
	level int
}

type Log interface {
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Write(p []byte) (n int, err error)
	Close() error
	SetMinLevel(level int)
	MinLevel() int
}

const (
	LogLevelFatal = iota
	LogLevelError
	LogLevelInfo
	LogLevelDebug
)

var (
	levelNames = [...]string{
		"FATAL", "ERROR", "INFO ", "DEBUG",
	}
	levelColors = [...]int{
		31, //red
		31, //red
		36, //blue
		37, //gray
	}
	ErrInvalidLogType = errors.New("Invalid log type")
)

const (
	LogTypeNull = iota
	LogTypeStdout
	LogTypeStdoutColor
	LogTypeSyslog
	LogTypeFile
)
