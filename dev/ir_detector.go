/*
IRDetector is a sensor used to detected infrared ray.

Connect to Pi:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 - out: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// IRDetector implements Detector interface
type IRDetector struct {
	pin rpio.Pin
}

// NewIRDetector ...
func NewIRDetector(pin uint8) *IRDetector {
	i := &IRDetector{
		pin: rpio.Pin(pin),
	}
	i.pin.Input()
	return i
}

// Detected ...
func (i *IRDetector) Detected() bool {
	return i.pin.Read() == rpio.Low
}
