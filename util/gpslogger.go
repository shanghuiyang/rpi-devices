package util

import (
	"fmt"
	"os"
	"time"

	"github.com/jakefau/rpi-devices/util/geo"
)

const (
	timeFormat = "2006-01-02T15:04:05"
)

// GPSLogger ...
type GPSLogger struct {
	f        *os.File
	chPoints chan *geo.Point
}

// NewGPSLogger ...
func NewGPSLogger() *GPSLogger {
	fname := time.Now().Format(timeFormat) + ".csv"
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil
	}
	f.WriteString("timestamp,lat,lon\n")
	t := &GPSLogger{
		f:        f,
		chPoints: make(chan *geo.Point, 32),
	}
	go t.start()
	return t
}

func (l *GPSLogger) start() {
	for pt := range l.chPoints {
		tm := time.Now().Format(timeFormat)
		line := fmt.Sprintf("%v,%.6f,%.6f\n", tm, pt.Lat, pt.Lon)
		l.f.WriteString(line)
	}
}

// AddPoint ...
func (l *GPSLogger) AddPoint(pt *geo.Point) {
	if pt == nil {
		return
	}
	l.chPoints <- pt
}

// Close ...
func (l *GPSLogger) Close() {
	l.f.Close()
	close(l.chPoints)
}
