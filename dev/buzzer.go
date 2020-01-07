/*
Package dev ...

Connect to Pi:
 - vcc: any v3.3 pin
 - gnd: and gnd pin
 - i/o: any data pin
*/
package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

// Buzzer ...
type Buzzer struct {
	pin rpio.Pin
}

// NewBuzzer ...
func NewBuzzer(pin int8) *Buzzer {
	b := &Buzzer{
		pin: rpio.Pin(pin),
	}
	b.pin.Output()
	return b
}

// On ...
func (b *Buzzer) On() {
	b.pin.High()
}

// Off ...
func (b *Buzzer) Off() {
	b.pin.Low()
}

// Beep beeps [n] times with an interval in [interval] millisecond
func (b *Buzzer) Beep(n int, interval int) {
	d := time.Duration(interval) * time.Millisecond
	for i := 0; i < n; i++ {
		b.pin.High()
		time.Sleep(d)
		b.pin.Low()
		time.Sleep(d)
	}
}
