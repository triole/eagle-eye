package conf

import (
	"eagle-eye/src/logging"
	"regexp"
	"time"
)

type Conf struct {
	Command      []string
	Interval     time.Duration
	Pause        time.Duration
	Folder       string
	Regex        *regexp.Regexp
	Spectate     bool
	KeepOutput   bool
	RunInitially bool
	Logging      logging.Logging
	LogInit      bool
	Verbose      bool
}

func Init() Conf {
	return Conf{}
}
