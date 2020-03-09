/*
Package dev ...

Connect to Pi:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 - out: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// Infrared ...
type Infrared struct {
	pin rpio.Pin
}

// NewInfrared ...
func NewInfrared(pin uint8) *Infrared {
	i := &Infrared{
		pin: rpio.Pin(pin),
	}
	i.pin.Input()
	return i
}

// Detected ...
func (i *Infrared) Detected() bool {
	return i.pin.Read() == rpio.Low
}
