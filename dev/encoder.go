/*
Package dev ...

Connect to Pi:
 - vcc: any 3.3v or v5 pin
 - gnd: any gnd pin
 - out: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// Encoder ...
type Encoder struct {
	pin rpio.Pin
}

// NewEncoder ...
func NewEncoder(pin uint8) *Encoder {
	e := &Encoder{
		pin: rpio.Pin(pin),
	}
	e.pin.Input()
	return e
}

// Detected ...
func (e *Encoder) Detected() bool {
	return e.pin.Read() == rpio.High
}

// Count1 ...
func (e *Encoder) Count1() int {
	if e.pin.Read() == rpio.High {
		return 1
	}
	return 0
}
