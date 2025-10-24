package conf

import (
	"regexp"
	"time"

	"github.com/triole/logseal"
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
	Lg           logseal.Logseal
	LogInit      bool
	Verbose      bool
}

func Init() Conf {
	return Conf{}
}
