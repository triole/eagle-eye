package main

import (
	"eagle-eye/logging"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type tSettings struct {
	Command    []string
	Interval   time.Duration
	Pause      time.Duration
	Folder     string
	Regex      *regexp.Regexp
	Spectate   bool
	KeepOutput bool
	Logging    logging.Logging
	LogInit    bool
	Verbose    bool
}

func main() {
	parseArgs()

	if CLI.PrintVars {
		printAvailableVars()
		os.Exit(0)
	}

	settings := tSettings{
		Command:    CLI.Command,
		Folder:     CLI.Folder,
		Interval:   time.Duration(CLI.Interval) * time.Second,
		Pause:      time.Duration(CLI.Pause) * time.Second,
		Regex:      regexp.MustCompile(CLI.Regex),
		Spectate:   CLI.Spectate,
		KeepOutput: CLI.KeepOutput,
		LogInit:    false,
	}
	if len(settings.Command) < 1 {
		settings.Spectate = true
	}

	settings.Logging = logging.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)

	mode := fmt.Sprintf("command on change: %q", settings.Command)
	if settings.Spectate {
		mode = "just spectate"
	}

	if CLI.RunInitially {
		color.Green("\nRun command initially %q, %+v", settings.Command)
		runCmd(settings.Command, settings.Pause, settings.Verbose)
	}

	settings.Logging.Info("Watch folder", logrus.Fields{
		"folder":       settings.Folder,
		"action":       mode,
		"logfile":      CLI.LogFile,
		"loglevel":     CLI.LogLevel,
		"lognocolours": CLI.LogNoColors,
		"logjson":      CLI.LogJSON,
	})

	watch(settings)
}
