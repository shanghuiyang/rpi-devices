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

// Collision ...
type Collision struct {
	pin rpio.Pin
}

// NewCollision ...
func NewCollision(pin uint8) *Collision {
	c := &Collision{
		pin: rpio.Pin(pin),
	}
	c.pin.Input()
	return c
}

// Collided ...
func (c *Collision) Collided() bool {
	return c.pin.Read() == rpio.Low
}
