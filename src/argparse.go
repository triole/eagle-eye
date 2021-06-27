package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	// BUILDTAGS are injected ld flags during build
	BUILDTAGS      string
	appName        = "eagle-eye"
	appDescription = "eagle-eye is a folder watcher that runs commands on change of files"
	appMainversion = "0.1"
)

var CLI struct {
	Command     []string `help:"command to run, flags always have to be in front" arg optional passthrough`
	Folder      string   `help:"folder to watch, default is current dir" optional short:f default:${curdir}`
	Regex       string   `help:"regex scheme, only consider files that match" optional short:r default:.*`
	Spectate    bool     `help:"spectate mode, do not run command, just watch and print file system changes" short:s`
	Interval    int32    `help:"watch interval, recheck for changes in seconds" default:1 short:i`
	KeepOutput  bool     `help:"keep output, do not clear screen when running command" short:k`
	LogFile     string   `help:"log file, if not given no logs will be written, this is default" short:l`
	VersionFlag bool     `help:"display version" short:V`
}

func parseArgs() {
	curdir, _ := os.Getwd()
	ctx := kong.Parse(&CLI,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:      true,
			Summary:      true,
			NoAppSummary: true,
			FlagsLast:    false,
		}),
		kong.Vars{
			"curdir": curdir,
		},
	)
	_ = ctx.Run()

	if CLI.VersionFlag == true {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
	// ctx.FatalIfErrorf(err)
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "Version: "+appMainversion+".", -1)
	fmt.Printf("%s\n", s)
}

func alnum(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile("[^a-z0-9_-]")
	return re.ReplaceAllString(s, "-")
}
