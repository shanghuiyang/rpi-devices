/*
Package dev ...

SW-420 is an sensor which is able to detect shaking.

Connect to Pi:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 - do : any data pin

*/
package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

// SW420 ...
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

// Shaked returns true if the sensor detects a shake,
// or return false
func (s *SW420) Shaked() bool {
	return s.pin.Read() == rpio.High
}

// KeepShaking returns true if the sensor detects the object keeps shaking in 100 millisecond,
// or returns false
func (s *SW420) KeepShaking() bool {
	states := map[bool]int{
		true:  0,
		false: 0,
	}
	for i := 0; i < 10; i++ {
		shaked := s.Shaked()
		states[shaked]++
		time.Sleep(10 * time.Millisecond)
	}
	return states[true] > states[false]
}
