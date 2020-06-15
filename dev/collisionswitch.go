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
	pin  rpio.Pin
	flag bool
}

// NewCollisionSwitch ...
func NewCollisionSwitch(pin uint8) *CollisionSwitch {
	c := &CollisionSwitch{
		pin:  rpio.Pin(pin),
		flag: true,
	}
	c.pin.Input()
	c.pin.PullUp()
	c.pin.Detect(rpio.FallEdge)
	return c
}

// Collided ...
func (c *CollisionSwitch) Collided() bool {
	collided := c.pin.EdgeDetected()
	if collided {
		if c.flag {
			c.flag = false
			return true
		}
		c.flag = true
	}
	return false
}
