/*
Package dev ...

SW-420 is an sensor which is able to detect shaking.

Connect to Pi:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 - do : any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// SW420 ...
type SW420 struct {
	pin rpio.Pin
}

// NewSW420 ...
func NewSW420(pin uint8) *SW420 {
	i := &SW420{
		pin: rpio.Pin(pin),
	}
	i.pin.Input()
	return i
}

// Shaked returns true if detect a shake, and return false if didn't detect any shakes
func (i *SW420) Shaked() bool {
	return i.pin.Read() == rpio.High
}
