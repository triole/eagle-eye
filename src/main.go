package main

import (
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
	Folder     string
	Regex      *regexp.Regexp
	Spectate   bool
	KeepOutput bool
	Logging    *logrus.Logger
	LogInit    bool
}

func main() {
	parseArgs()

	if CLI.PrintVars == true {
		printAvailableVars()
		os.Exit(0)
	}

	settings := tSettings{
		Command:    CLI.Command,
		Folder:     CLI.Folder,
		Interval:   time.Duration(CLI.Interval) * time.Second,
		Regex:      regexp.MustCompile(CLI.Regex),
		Spectate:   CLI.Spectate,
		KeepOutput: CLI.KeepOutput,
		LogInit:    false,
	}
	if len(settings.Command) < 1 {
		settings.Spectate = true
	}

	if CLI.LogFile != "" {
		settings.Logging, settings.LogInit = initLogging(CLI.LogFile)
	}

	mode := fmt.Sprintf("command on change: %q", settings.Command)
	if settings.Spectate == true {
		mode = "just spectate"
	}

	if CLI.RunInitially == true {
		color.Green("\nRun command initially %q, %+v", settings.Command)
		runCmd(settings.Command)
	}

	color.Green("\nWatch folder %q, %s", settings.Folder, mode)
	watch(settings)
}
