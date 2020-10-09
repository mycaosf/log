package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

func (p *logFile) logDefault(format string, args ...interface{}) {
	w := p.flog
	fmt.Fprintf(w, "%s %s ", time.Now().Local().Format(timeFormatDefault), levelNames[p.level])
	fmt.Fprintf(w, format, args...)
	fmt.Fprintf(w, "\n")
}

func (p *logFile) Fatalf(format string, args ...interface{}) {
	if p.level >= LogLevelFatal {
		p.log(format, args...)
		os.Exit(1)
	}
}

func (p *logFile) Errorf(format string, args ...interface{}) {
	if p.level >= LogLevelError {
		p.log(format, args...)
	}
}

func (p *logFile) Infof(format string, args ...interface{}) {
	if p.level >= LogLevelInfo {
		p.log(format, args...)
	}
}

func (p *logFile) Debugf(format string, args ...interface{}) {
	if p.level >= LogLevelDebug {
		p.log(format, args...)
	}
}

func (p *logFile) Write(b []byte) (int, error) {
	return p.flog.Write(b)
}

func (p *logFile) Close() error {
	if p.flog != nil {
		if p.flog != os.Stdout && p.flog != os.Stderr {
			p.flog.Close()
		}
		p.flog = nil
	}

	return nil
}

func NewFile(name string, minLevel int) (l Log, err error) {
	var flog *os.File
	if flog, err = os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err == nil {
		l = newFile(minLevel, flog)
	}

	return
}

func NewStdout(minLevel int) Log {
	return newFile(minLevel, os.Stdout)
}

func newFile(minLevel int, file io.WriteCloser) Log {
	ret := &logFile{
		logBase: logBase{level: minLevel},
		flog:    file,
	}
	ret.log = ret.logDefault

	return ret
}

type logFile struct {
	logBase
	flog io.WriteCloser
	log  func(format string, args ...interface{})
}

const (
	timeFormatDefault = "2006-01-02 15:04:05"
)
