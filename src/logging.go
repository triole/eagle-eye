package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func initLogging(logFile string) (log *logrus.Logger, b bool) {
	log = logrus.New()
	b = true
	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "date",
			logrus.FieldKeyMsg:  "operation",
		},
		TimestampFormat: "2006-01-02 15:04:05.000 MST",
		PrettyPrint:     false,
	})

	logFileOpen, err := os.OpenFile(
		logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666,
	)
	if err != nil {
		log.Info("Failed to log to file, use default stderr")
		b = false
	}

	logrus.SetOutput(logFileOpen)
	return
}
