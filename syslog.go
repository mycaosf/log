package log

import (
	"fmt"
	syslog "github.com/racksec/srslog"
	"os"
)

func (p *logSyslog) Fatalf(format string, args ...interface{}) {
	if p.level >= LogLevelFatal {
		str := fmt.Sprintf(format, args...)
		p.slog.Crit(str)
		os.Exit(1)
	}
}

func (p *logSyslog) Errorf(format string, args ...interface{}) {
	if p.level >= LogLevelError {
		str := fmt.Sprintf(format, args...)
		p.slog.Err(str)
	}
}

func (p *logSyslog) Infof(format string, args ...interface{}) {
	if p.level >= LogLevelInfo {
		str := fmt.Sprintf(format, args...)
		p.slog.Info(str)
	}
}

func (p *logSyslog) Debugf(format string, args ...interface{}) {
	if p.level >= LogLevelDebug {
		str := fmt.Sprintf(format, args...)
		p.slog.Debug(str)
	}
}

func (p *logSyslog) Write(b []byte) (int, error) {
	return p.slog.Write(b)
}

func (p *logSyslog) Close() error {
	if p.slog != nil {
		p.slog.Close()
		p.slog = nil
	}

	return nil
}

func NewSyslog(tag string, minLevel int) (l Log, err error) {
	var slog *syslog.Writer
	if slog, err = syslog.Dial("", "", syslog.LOG_INFO, tag); err == nil {
		l = &logSyslog{
			logBase: logBase{
				level: minLevel,
			},
			slog: slog,
		}
	}

	return
}

type logSyslog struct {
	logBase
	slog *syslog.Writer
}
