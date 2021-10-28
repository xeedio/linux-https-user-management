package humcommon

import (
	"io"
	"log/syslog"

	"github.com/sirupsen/logrus"

	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

//nolint
func main() {
}

var logger *logrus.Logger

func initLogger() {
	logger = logrus.New()
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.WarnLevel)

	if err == nil {
		logger.Hooks.Add(hook)
		logger.SetOutput(io.Discard)
	}
}

func Log() *logrus.Logger {
	return logger
}
