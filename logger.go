package humcommon

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
)

var logPrefix = "LHUM "
var logger *log.Logger

func SetLogPrefix(p string) {
	logPrefix = p
}

func initLogger() error {
	var err error
	if AppConfig.Logger == "syslog" {
		logger, err = syslog.NewLogger(syslog.LOG_INFO, log.Ltime)
		if err != nil {
			logger = log.New(os.Stderr, logPrefix, log.Ltime)
			logger.Printf("syslog.Open() err: %v", err)
		}
	} else {
		logger = log.New(os.Stdout, logPrefix, log.Ltime)
	}

	return nil
}

func LogDebug(module string, data ...interface{}) {
	if !AppConfig.Debug {
		return
	}
	module = "DEBUG-" + module
	LogInfo(module, data)
}

func LogInfo(module string, data ...interface{}) {
	s := fmt.Sprintf("[%s] %s", module, fmt.Sprint(data...))
	logger.Println(s)
}

func LogFatal(module string, err error) {
	s := fmt.Sprintf("[%s] Fatal: %v", module, err)
	logger.Fatal(s)
}
