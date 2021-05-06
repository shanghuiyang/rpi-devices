/*
Package dev ...

Connect to Pi:
 - S: 	any data pin
 - M:   5v
 - -:   any ground

*/
package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// Collision ...
type ShockSensor struct {
	pin rpio.Pin
}

func NewShockSensor(pin uint8) *ShockSensor{
	s := &ShockSensor{
		pin: rpio.Pin(pin),
	}
	s.pin.Input()
	return s
}

// Shocked ...
func (s *ShockSensor) Shock() bool {
	return s.pin.Read() == rpio.Low
}
