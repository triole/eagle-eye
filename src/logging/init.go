package logging

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	logLevels = map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"error": logrus.ErrorLevel,
	}
)

// Logging holds the logging module
type Logging struct {
	Logrus    *logrus.Logger
	LogToFile bool
}

// Init method, does what it says
func Init(loglevel, logFile string, nocolours, JSONLog bool) (lg Logging) {
	timeStampFormat := "2006-01-02 15:04:05.000 MST"
	lg.Logrus = logrus.New()

	if JSONLog {
		lg.Logrus.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "date",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "msg",
			},
			TimestampFormat:   timeStampFormat,
			PrettyPrint:       false,
			DisableHTMLEscape: false,
		})
	} else {
		form := &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: timeStampFormat,
			DisableQuote:    true,
			PadLevelText:    true,
			ForceColors:     true,
		}
		if nocolours {
			form.ForceColors = false
			form.DisableColors = true
		}
		lg.Logrus.SetFormatter(form)
	}

	if logFile != "/dev/stdout" {
		lg.LogToFile = true
		openLogFile, err := os.OpenFile(
			logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644,
		)
		lg.IfErrFatal(
			"Can not open log file",
			F{
				"logfile": logFile,
				"error":   err.Error(),
			},
		)

		lg.Logrus.SetOutput(openLogFile)
	}
	lg.setLevel(loglevel)

	return lg
}

func (lg *Logging) setLevel(level string) {
	if val, ok := logLevels[strings.ToLower(level)]; ok {
		lg.Logrus.SetLevel(val)
	} else {
		lg.setLevel("info")
	}
}
