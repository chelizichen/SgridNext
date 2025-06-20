package schedule

import (
	"time"
)

func wrapFunc(t *time.Ticker, cb func()) {
	for range t.C {
		cb()
	}
}

func LoadTick() {
	go wrapFunc(time.NewTicker(30*time.Second), runRestartCallback)
}
