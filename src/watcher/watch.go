package watcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/triole/logseal"
)

var (
	lastRun     time.Time
	debounceMu  sync.Mutex    // Mutex to protect debounce state
	debounceEnd time.Time     // End time of debounce period
	debounceCh  chan struct{} // Channel to signal debounce complete
)

// waitForDebounce blocks until the debounce period has elapsed
func waitForDebounce(debounce time.Duration) {
	if debounce == 0 {
		return
	}
	debounceMu.Lock()
	debounceEnd = time.Now().Add(debounce)
	debounceMu.Unlock()
	<-time.After(debounce)
}

// shouldDebounce checks if we should wait before running the command
// Returns true if we should skip execution (still in debounce period)
func shouldDebounce() bool {
	debounceMu.Lock()
	defer debounceMu.Unlock()
	if time.Now().Before(debounceEnd) {
		return true
	}
	return false
}

type Event struct {
	Time  time.Time
	Event watcher.Event
}

type EventChan chan Event

type tVarMapEntry struct {
	Val  interface{}
	Desc string
}

func (w Watcher) Run() {
	var err error
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
				w.Conf.Lg.Fatal("An error occured", logseal.F{
					"error": err,
				})
			case <-w.Watcher.Closed:
				return
			}
		}
	}()

	err = w.Watcher.AddRecursive(w.Conf.Folder)
	w.Conf.Lg.IfErrFatal("Unable to add folders to watch list", logseal.F{
		"error": err,
	})

	go func() {
		w.Watcher.Wait()
	}()

	err = w.Watcher.Start(w.Conf.Interval)
	w.Conf.Lg.IfErrFatal("Can not start watcher", logseal.F{
		"error": err,
	})
}

func (w Watcher) runChannelWatcher(chin EventChan) {
	// cache config values for better performance
	interval := w.Conf.Interval
	spectate := w.Conf.Spectate
	keepOutput := w.Conf.KeepOutput
	pause := w.Conf.Pause
	verbose := w.Conf.Verbose
	command := w.Conf.Command
	debounce := w.Conf.Debounce

	for ev := range chin {
		if ev.Event.Op > 0 {
			if !spectate && w.calcDiff(lastRun, interval) {
				// Wait for debounce period to complete if there's pending activity
				waitForDebounce(debounce)

				// Check if we're still in debounce period (another event may have extended it)
				if shouldDebounce() {
					w.printEvent(ev.Event)
					continue
				}

				if !keepOutput {
					fmt.Print("\033[2J")
					fmt.Print("\033[H")
				}

				lastRun = time.Now()

				cmdArr := w.iterTemplate(command, w.makeVarMap(ev.Event))
				w.Conf.Lg.Info("Run", logseal.F{
					"cmds": cmdArr,
				})
				w.runCmd(cmdArr, pause, verbose)
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
		w.Conf.Lg.Info(
			"Event",
			logseal.F{
				"event": event.Op.String(),
				"path":  fmt.Sprintf("%s", event.Path),
				"type":  t,
			},
		)
	} else {
		w.Conf.Lg.Debug("Event", logseal.F{
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
