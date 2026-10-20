package watcher

import (
	"eagle-eye/src/conf"
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
}

func Init(conf conf.Conf) (w Watcher) {
	w = Watcher{
		Conf:    conf,
		Watcher: watcher.New(),
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
