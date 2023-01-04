package watcher

import (
	"eagle-eye/src/conf"
	"eagle-eye/src/logging"

	"github.com/radovskyb/watcher"
)

type Watcher struct {
	Conf    conf.Conf
	Watcher *watcher.Watcher
}

func Init(conf conf.Conf) (w Watcher) {
	w = Watcher{
		Conf:    conf,
		Watcher: watcher.New(),
	}
	w.Watcher.AddFilterHook(watcher.RegexFilterHook(conf.Regex, false))

	if w.Conf.RunInitially {
		conf.Logging.Info("Run initially", logging.F{
			"cmds": conf.Command,
		})
		w.runCmd(conf.Command, conf.Pause, conf.Verbose)
	}

	return
}
