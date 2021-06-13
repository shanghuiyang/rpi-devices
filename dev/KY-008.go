/*Package dev ...

Connect to Pi:
- positive: 	    	any data pin
- negative(marked): 	any gnd pin
*/

package dev

import (
	"github.com/stianeikeland/go-rpio"
	"time"
)

const (
	logTagLaser = "laser"
)

// Laser ...
type Laser struct {
	pin rpio.Pin
}

// NewLaser...
func NewLaser(pin uint8) *Laser {
	l := &Laser{
		pin: rpio.Pin(pin),
	}
	l.pin.Output()
	return l
}

// On ...
func (l *Laser) On() {
	l.pin.High()
}

// Off ...
func (l *Laser) Off() {
	l.pin.Low()
}

// Blink is let led blink n time, interval Millisecond each time
func (l *Laser) Blink(n int, interval int) {
	d := time.Duration(interval) * time.Millisecond
	for i := 0; i < n; i++ {
		l.On()
		time.Sleep(d)
		l.Off()
		time.Sleep(d)
	}
}
