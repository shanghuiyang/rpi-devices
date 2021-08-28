/*
SW-420 is an sensor used to detect shaking.

Connect to Pi:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 - do : any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// SW420 implements Detector interface
type SW420 struct {
	pin rpio.Pin
}

// NewSW420 ...
func NewSW420(pin uint8) *SW420 {
	s := &SW420{
		pin: rpio.Pin(pin),
	}
	s.pin.Input()
	return s
}

// Detected returns true if the sensor detects shaking,
// or return false
func (s *SW420) Detected() bool {
	return s.pin.Read() == rpio.High
}
