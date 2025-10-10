package schedule

import (
	"time"

	"sgridnext.com/src/probe"
)

func wrapFunc(t *time.Ticker, cb func()) {
	for range t.C {
		cb()
	}
}



func loadProbe() {
	go wrapFunc(time.NewTicker(3*time.Hour), probe.RunProbeTask)
}


func LoadSchedule(){
	loadProbe()
}