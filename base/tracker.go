package base

import (
	"fmt"
	"os"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05"
)

// Tracker ...
type Tracker struct {
	f  *os.File
	ch chan *Point
}

// NewTracker ...
func NewTracker() *Tracker {
	fname := time.Now().Format(timeFormat) + ".csv"
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil
	}
	f.WriteString("timestamp,lat,lon\n")
	t := &Tracker{
		f:  f,
		ch: make(chan *Point, 32),
	}
	go t.start()
	return t
}

func (t *Tracker) start() {
	for pt := range t.ch {
		tm := time.Now().Format(timeFormat)
		line := fmt.Sprintf("%v,%.6f,%.6f\n", tm, pt.Lat, pt.Lon)
		t.f.WriteString(line)
	}
}

// AddPoint ...
func (t *Tracker) AddPoint(pt *Point) {
	t.ch <- pt
}

// Close ...
func (t *Tracker) Close() {
	t.f.Close()
	close(t.ch)
}
