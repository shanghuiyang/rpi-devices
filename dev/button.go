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

// Button ...
type Button struct {
	pin rpio.Pin
}

// NewButton ...
func NewButton(pin uint8) *Button {
	b := &Button{
		pin: rpio.Pin(pin),
	}
	b.pin.Input()
	return b
}

// Pressed ...
func (i *Button) Pressed() bool {
	return i.pin.Read() == rpio.High
}
