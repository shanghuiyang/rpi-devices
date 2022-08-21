/*
Fan is a fan module using 3.3v power source.

Connect to Raspberry Pi:
  - vcc(red line)  : any data pin(~3.3v)
  - gnd(black line): any gnd pin
*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// Fan ...
type Fan struct {
	pin rpio.Pin
}

// NewFan ...
func NewFan(pin uint8) *Fan {
	f := &Fan{
		pin: rpio.Pin(pin),
	}
	f.pin.Output()
	f.pin.Low()
	return f
}

// On ...
func (f *Fan) On() {
	f.pin.High()
}

// Off ...
func (f *Fan) Off() {
	f.pin.Low()
}
