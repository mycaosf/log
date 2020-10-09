package log

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"io"
	"time"
)

func NewStdoutColor(minLevel int) Log {
	w := colorable.NewColorableStdout()
	ret := &logFileColor{
		logFile: logFile{
			logBase: logBase{level: minLevel},
			flog:    &dummyCloser{w: w},
		},
	}
	ret.log = ret.logColor

	return ret
}

func (p *logFileColor) logColor(format string, args ...interface{}) {
	w := p.flog
	fmt.Fprintf(w, "%s \x1b[%dm%s ", time.Now().Local().Format(timeFormatDefault), levelColors[p.level], levelNames[p.level])
	fmt.Fprintf(w, format, args...)
	fmt.Fprintf(w, "\x1b[0m\n")
}

type logFileColor struct {
	logFile
}

func (p *dummyCloser) Write(b []byte) (int, error) {
	return p.w.Write(b)
}

func (p *dummyCloser) Close() error {
	return nil
}

type dummyCloser struct {
	w io.Writer
}
