package main

import (
	"fmt"
	"log"
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

func watch(settings tSettings) {
	w := watcher.New()
	w.AddFilterHook(watcher.RegexFilterHook(settings.Regex, false))

	chin := make(EventChan)
	go ticker(chin)
	go runChannelWatcher(settings, chin)

	go func() {
		for {
			select {
			case ev := <-w.Event:
				event := Event{
					Time:  time.Now(),
					Event: ev,
				}
				chin <- event
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(settings.Folder); err != nil {
		log.Fatalln(err)
	}

	go func() {
		w.Wait()
	}()

	if err := w.Start(settings.Interval); err != nil {
		log.Fatalln(err)
	}
}

func runChannelWatcher(settings tSettings, chin EventChan) {
	current := time.Now()
	last := time.Now().Add(-time.Second * (settings.Interval * 2))
	diff := current.Sub(last) > settings.Interval-settings.Interval/4
	var lastDiff bool
	for ev := range chin {
		if ev.Event.Op > 0 {
			if calcDiff(lastRun, settings.Interval) && !settings.Spectate {
				if !settings.KeepOutput {
					fmt.Print("\033[2J")
					fmt.Print("\033[H")
				}
				now := time.Now()
				lastRun = now
				lastDiff = diff
				last = current
				current = ev.Time
				if lastDiff && calcDiff(last, settings.Interval) {
					lastRun = time.Now()
					cmdArr := iterTemplate(
						settings.Command, makeVarMap(ev.Event),
					)
					settings.Logging.Info("Run", logrus.Fields{
						"cmds": cmdArr,
					})
					runCmd(cmdArr, settings.Pause, settings.Verbose)
				}
			} else {
				printEvent(ev.Event, settings)
			}
		}
	}
}

func calcDiff(lastRun time.Time, interval time.Duration) bool {
	return time.Since(lastRun) > interval
}

func ticker(chin EventChan) {
	for range time.Tick(time.Duration(1) * time.Second) {
		event := Event{
			Time: time.Now(),
		}
		chin <- event
	}
}

func printEvent(event watcher.Event, settings tSettings) {
	t := "FILE"
	if event.IsDir() {
		t = "FOLDER"
	}
	if settings.Spectate {
		settings.Logging.Info("Event", logrus.Fields{
			"event": event.Op.String(),
			"path":  fmt.Sprintf(event.Path),
			"type":  t,
		})
	} else {
		settings.Logging.Debug("Event", logrus.Fields{
			"event": event.Op.String(),
		})
	}
}

func makeVarMap(ev watcher.Event) (varMap map[string]tVarMapEntry) {
	varMap = make(map[string]tVarMapEntry)
	varMap["file"] = tVarMapEntry{
		ev.Path, "file that triggered the event",
	}
	varMap["folder"] = tVarMapEntry{
		find(`.*/`, ev.Path), "folder of the file that triggered the event",
	}
	return
}
