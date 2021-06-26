package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/fatih/color"
	"github.com/radovskyb/watcher"
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

func watch(interval time.Duration) {
	w := watcher.New()

	r := regexp.MustCompile(CLI.Regex)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	chin := make(EventChan)
	go ticker(chin)
	go runChannelWatcher(interval, chin)

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

	if err := w.AddRecursive(CLI.Folder); err != nil {
		log.Fatalln(err)
	}

	go func() {
		w.Wait()
	}()

	if err := w.Start(interval); err != nil {
		log.Fatalln(err)
	}
}

func runChannelWatcher(interval time.Duration, chin EventChan) {
	current := time.Now()
	last := time.Now()
	diff := current.Sub(last) > interval-interval/4
	var lastDiff bool
	for ev := range chin {
		if ev.Event.Op > 0 {
			if time.Now().Sub(lastRun) > interval+interval/4 {
				fmt.Print("\033[2J")
				fmt.Print("\033[H")
				now := time.Now()
				color.Green("\n%s\n", now)
				lastRun = now
			}
			printEvent(ev.Event)
		}
		if CLI.Spectate == false {
			lastDiff = diff
			last = current
			current = ev.Time
			diff = current.Sub(last) > interval-interval/4
			if lastDiff == false && diff == true {
				lastRun = time.Now()
				runCmd(CLI.Command, true)
			}
		}
	}
}

func ticker(chin EventChan) {
	for _ = range time.Tick(time.Duration(1) * time.Second) {
		event := Event{
			Time: time.Now(),
		}
		chin <- event
	}
}

func printEvent(event watcher.Event) {
	t := "FILE"
	if event.IsDir() == true {
		t = "FOLDER"
	}
	if event.Path == event.OldPath {
		fmt.Printf("%s\t%s\t%s\n", t, event.Op, event.Path)
	} else {
		fmt.Printf("%s\t%s\t%s %s\n", t, event.Op, event.Path, event.OldPath)
	}
}
