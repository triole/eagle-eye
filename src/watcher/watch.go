package watcher

import (
	"fmt"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/sirupsen/logrus"
)

var (
	lastRun time.Time
)

// Event holds an event and its time
type Event struct {
	Time  time.Time
	Event watcher.Event
}

// EventChan holds what is pushed into the processing channels
type EventChan chan Event

type tVarMapEntry struct {
	Val  interface{}
	Desc string
}

func (w Watcher) Run() {

	chin := make(EventChan)
	go w.ticker(chin)
	go w.runChannelWatcher(chin)

	go func() {
		for {
			select {
			case ev := <-w.Watcher.Event:
				event := Event{
					Time:  time.Now(),
					Event: ev,
				}
				chin <- event
			case err := <-w.Watcher.Error:
				w.Conf.Logging.Fatal("An error occured", logrus.Fields{
					"error": err,
				})
			case <-w.Watcher.Closed:
				return
			}
		}
	}()

	if err := w.Watcher.AddRecursive(w.Conf.Folder); err != nil {
		w.Conf.Logging.Fatal("Unable to add folders to watch list", logrus.Fields{
			"error": err,
		})
	}

	go func() {
		w.Watcher.Wait()
	}()

	if err := w.Watcher.Start(w.Conf.Interval); err != nil {
		w.Conf.Logging.Fatal("Can not start watcher", logrus.Fields{
			"error": err,
		})
	}
}

func (w Watcher) runChannelWatcher(chin EventChan) {
	current := time.Now()
	last := time.Now().Add(-time.Second * (w.Conf.Interval * 2))
	diff := current.Sub(last) > w.Conf.Interval-w.Conf.Interval/4
	var lastDiff bool
	for ev := range chin {
		if ev.Event.Op > 0 {
			if w.calcDiff(lastRun, w.Conf.Interval) && !w.Conf.Spectate {
				if !w.Conf.KeepOutput {
					fmt.Print("\033[2J")
					fmt.Print("\033[H")
				}
				now := time.Now()
				lastRun = now
				lastDiff = diff
				last = current
				current = ev.Time
				if lastDiff && w.calcDiff(last, w.Conf.Interval) {
					lastRun = time.Now()
					cmdArr := w.iterTemplate(
						w.Conf.Command, w.makeVarMap(ev.Event),
					)
					w.Conf.Logging.Info("Run", logrus.Fields{
						"cmds": cmdArr,
					})
					w.runCmd(cmdArr, w.Conf.Pause, w.Conf.Verbose)
				}
			} else {
				w.printEvent(ev.Event)
			}
		}
	}
}

func (w Watcher) calcDiff(lastRun time.Time, interval time.Duration) bool {
	return time.Since(lastRun) > interval
}

func (w Watcher) ticker(chin EventChan) {
	for range time.Tick(time.Duration(1) * time.Second) {
		event := Event{
			Time: time.Now(),
		}
		chin <- event
	}
}

func (w Watcher) printEvent(event watcher.Event) {
	t := "FILE"
	if event.IsDir() {
		t = "FOLDER"
	}
	if w.Conf.Spectate {
		w.Conf.Logging.Info("Event", logrus.Fields{
			"event": event.Op.String(),
			"path":  fmt.Sprintf(event.Path),
			"type":  t,
		})
	} else {
		w.Conf.Logging.Debug("Event", logrus.Fields{
			"event": event.Op.String(),
		})
	}
}

func (w Watcher) makeVarMap(ev watcher.Event) (varMap map[string]tVarMapEntry) {
	varMap = make(map[string]tVarMapEntry)
	varMap["file"] = tVarMapEntry{
		ev.Path, "file that triggered the event",
	}
	varMap["folder"] = tVarMapEntry{
		w.rxfind(`.*/`, ev.Path), "folder of the file that triggered the event",
	}
	return
}
