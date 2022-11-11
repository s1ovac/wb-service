package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Init() *logrus.Logger {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)

	log.SetLevel(logrus.InfoLevel)
	return log
}
