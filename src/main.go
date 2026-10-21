package main

import (
	"eagle-eye/src/conf"
	"eagle-eye/src/watcher"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/triole/logseal"
)

func main() {
	parseArgs()

	conf := conf.Init()
	conf.Lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)

	conf.Command = CLI.Command
	pth, err := filepath.Abs(CLI.Folder)
	conf.Lg.IfErrFatal(
		"can not determine absolute path: %s",
		logseal.F{"path": pth, "error": err},
	)
	conf.Folder = pth
	conf.Interval = time.Duration(CLI.Interval) * time.Second
	conf.Pause = time.Duration(CLI.Pause) * time.Second
	conf.Debounce = time.Duration(CLI.Debounce) * time.Millisecond
	conf.Regex = regexp.MustCompile(CLI.Regex)
	conf.Spectate = CLI.Spectate
	conf.KeepOutput = CLI.KeepOutput
	conf.RunInitially = CLI.RunInitially

	if len(conf.Command) < 1 {
		conf.Spectate = true
	}

	mode := fmt.Sprintf("run on change: %q", conf.Command)
	if conf.Spectate {
		mode = "spectate"
	}

	conf.Lg.Info("Watch folder", logseal.F{
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
