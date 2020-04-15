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
	b.pin.PullDown()
	b.pin.Detect(rpio.RiseEdge)
	return b
}

// Pressed ...
func (b *Button) Pressed() bool {
	return b.pin.EdgeDetected()
}
