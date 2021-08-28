/*
CollisionDetecter is a module used to detect whether a collision is happening.

Connect to Pi:
 - vcc: any 3.3v or v5 pin
 - gnd: any gnd pin
 - out: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// CollisionDetector implements Detector interface
type CollisionDetector struct {
	pin rpio.Pin
}

// NewCollisionDetector ...
func NewCollisionDetector(pin uint8) *CollisionDetector {
	c := &CollisionDetector{
		pin: rpio.Pin(pin),
	}
	c.pin.Input()
	return c
}

// Collided ...
func (c *CollisionDetector) Detected() bool {
	return c.pin.Read() == rpio.Low
}
