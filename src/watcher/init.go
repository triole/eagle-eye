package watcher

import (
	"eagle-eye/src/conf"
	"regexp"
	"sync"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/triole/logseal"
)

type Watcher struct {
	Conf          conf.Conf
	Watcher       *watcher.Watcher
	lastRunTime   time.Time
	debounceMu    sync.Mutex
	debounceTimer *time.Timer
	regexMu       sync.RWMutex
	regexCache    map[string]*regexp.Regexp
}

func Init(conf conf.Conf) (w Watcher) {
	w = Watcher{
		Conf:       conf,
		Watcher:    watcher.New(),
		regexCache: make(map[string]*regexp.Regexp),
	}
	w.Watcher.AddFilterHook(watcher.RegexFilterHook(conf.Regex, false))

	if w.Conf.RunInitially {
		conf.Lg.Info("Run initially", logseal.F{
			"cmds": conf.Command,
		})
		w.runCmd(conf.Command, conf.Pause, conf.Verbose)
	}

	return
}
