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

// CollisionSwitch ...
type CollisionSwitch struct {
	pin rpio.Pin
}

// NewCollisionSwitch ...
func NewCollisionSwitch(pin uint8) *CollisionSwitch {
	c := &CollisionSwitch{
		pin: rpio.Pin(pin),
	}
	c.pin.Input()
	return c
}

// Collided ...
func (c *CollisionSwitch) Collided() bool {
	return c.pin.Read() == rpio.Low
}
