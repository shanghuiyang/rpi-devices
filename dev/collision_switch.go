/*
CollisionSwitch is used to detect a collision with an object or to detect the limit of travel.

Connect to Raspberry Pi:
 - vcc: any 3.3v or v5 pin
 - gnd: any gnd pin
 - out: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// CollisionSwitch implements Detector interface
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

// Deteched ...
func (c *CollisionSwitch) Detected() bool {
	return c.pin.Read() == rpio.Low
}
