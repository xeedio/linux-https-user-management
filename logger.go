package humcommon

import (
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

	if err == nil {
		logger.Hooks.Add(hook)
	}
}

func Log() *logrus.Logger {
	return logger
}
