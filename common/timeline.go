package common

import (
	"fmt"
	"time"
)

var (
	TimelineEnabled bool
	pipe            chan string
	done            chan bool
)

type timevent struct {
	when int64
	what string
}

func timeliner(events chan string, done chan bool) {
	birth := time.Now()
	history := make([]*timevent, 0, 100)
	for {
		event, ok := <-events
		if !ok {
			break
		}
		history = append(history, &timevent{time.Since(birth).Milliseconds(), event})
	}
	death := time.Since(birth).Milliseconds()
	if TimelineEnabled && death > 0 {
		history = append(history, &timevent{death, "Now."})
		Log("----  rcc timeline  ----")
		Log(" #  percent  millis  event")
		for at, event := range history {
			permille := event.when * 1000 / death
			percent := float64(permille) / 10.0
			Log("%2d:  %5.1f%%  %6d  %s", at+1, percent, event.when, event.what)
		}
		Log("----  rcc timeline  ----")
	}
	close(done)
}

func init() {
	pipe = make(chan string)
	done = make(chan bool)
	go timeliner(pipe, done)
}

func Timeline(form string, details ...interface{}) {
	pipe <- fmt.Sprintf(form, details...)
}

func EndOfTimeline() {
	close(pipe)
	<-done
}
