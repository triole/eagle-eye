package main

import (
	"eagle-eye/src/conf"
	"eagle-eye/src/logging"
	"eagle-eye/src/watcher"
	"fmt"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	parseArgs()

	conf := conf.Init()
	conf.Command = CLI.Command
	conf.Folder = CLI.Folder
	conf.Interval = time.Duration(CLI.Interval) * time.Second
	conf.Pause = time.Duration(CLI.Pause) * time.Second
	conf.Regex = regexp.MustCompile(CLI.Regex)
	conf.Spectate = CLI.Spectate
	conf.KeepOutput = CLI.KeepOutput
	conf.RunInitially = CLI.RunInitially

	if len(conf.Command) < 1 {
		conf.Spectate = true
	}

	conf.Logging = logging.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)

	mode := fmt.Sprintf("run on change: %q", conf.Command)
	if conf.Spectate {
		mode = "just spectate"
	}

	conf.Logging.Info("Watch folder", logrus.Fields{
		"folder":         conf.Folder,
		"action":         mode,
		"log-file":       CLI.LogFile,
		"log-level":      CLI.LogLevel,
		"log-no-colours": CLI.LogNoColors,
		"log-json":       CLI.LogJSON,
	})

	w := watcher.Init(conf)

	if CLI.PrintVars {
		w.PrintAvailableVars()
	} else {
		w.Run()
	}
}
