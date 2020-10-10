# log
Simple log support stdout, file, syslog

Examples
```go
package main

import (
	"github.com/mycaosf/log"
)

func main() {
	//log, err := log.New(log.LogTypeSyslog, log.LogLevelDebug, "simple")
	log, err := log.New(log.LogTypeStdoutColor, log.LogLevelDebug, "simple")
	if err == nil {
		defer log.Close()
		log.Debugf("debug message: %d", 1)
		log.Infof("info message: %d", 1)
		log.Errorf("error message: %d", 1)
		log.Fatalf("fatalf message: %d", 1)
	}
}
```
