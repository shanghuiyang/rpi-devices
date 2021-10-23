/*
ButtonImp is a button module used to detect whether a button is pressed.

Connect to Raspberry Pi for a 3-pin buttom mobule:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 - out: any data pin

Connect to Raspberry Pi for a 2-pin buttom mobule:
 - port-1: any 3.3v pin
 - port-2: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// ButtonImp implements Button interface
type ButtonImp struct {
	pin rpio.Pin
}

// NewButtonImp ...
func NewButtonImp(pin uint8) *ButtonImp {
	b := &ButtonImp{
		pin: rpio.Pin(pin),
	}
	b.pin.Input()
	return b
}

// Pressed ...
func (b *ButtonImp) Pressed() bool {
	return b.pin.Read() == rpio.High
}
