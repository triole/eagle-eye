package watcher

import (
	"eagle-eye/src/conf"

	"github.com/radovskyb/watcher"
	"github.com/sirupsen/logrus"
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
		conf.Logging.Info("Run initially", logrus.Fields{
			"cmds": conf.Command,
		})
		w.runCmd(conf.Command, conf.Pause, conf.Verbose)
	}

	return
}
